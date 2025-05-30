// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package wsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	common "github.com/e154/smart-home/internal/system/gate/common"

	"github.com/gorilla/websocket"
)

// ConnectionStatus is an enumeration type which represents the status of WebSocket connection.
type ConnectionStatus int

var (
	upgrader = websocket.Upgrader{}
)

const (
	// Idle state means it is opened but not working now.
	// The default value for Connection is Idle, so it is ok to use zero-value(int: 0) for Idle status.
	Idle ConnectionStatus = iota
	Busy
	Closed
)

// Connection manages a single websocket connection from the peer.
// wsp supports multiple connections from a single peer at the same time.
type Connection struct {
	pool      *Pool
	ws        *websocket.Conn
	status    ConnectionStatus
	idleSince time.Time
	queue     chan Message
}

// NewConnection returns a new Connection.
func NewConnection(pool *Pool, ws *websocket.Conn) *Connection {
	c := &Connection{
		pool:   pool,
		ws:     ws,
		status: Idle,
		queue:  make(chan Message),
	}

	// Mark that this connection is ready to use for relay
	c.Release()

	return c
}

func (c *Connection) WritePump() {
	var data []byte
	var messageType int
	var err error
	for c.status != Closed {
		messageType, data, err = c.ws.ReadMessage()
		if messageType == -1 || err != nil {
			c.status = Closed
			close(c.queue)
			return
		}
		msg := Message{
			Type:  messageType,
			Value: data,
		}
		c.queue <- msg
	}
}

func (c *Connection) proxyWs(w http.ResponseWriter, r *http.Request) (err error) {
	defer c.Release()

	// Only pass those headers to the upgrader.
	upgradeHeader := http.Header{}
	if hdr := r.Header.Get("Sec-Websocket-Protocol"); hdr != "" {
		upgradeHeader.Set("Sec-Websocket-Protocol", hdr)
	}
	if hdr := r.Header.Get("Set-Cookie"); hdr != "" {
		upgradeHeader.Set("Set-Cookie", hdr)
	}

	query := r.URL.Query()
	accessToken := query.Get("access_token")
	if accessToken == "" {
		accessToken = "NIL"
	}
	if err = c.ws.WriteMessage(websocket.TextMessage, []byte("WS:"+accessToken)); err != nil {
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	connPub, err := upgrader.Upgrade(w, r, upgradeHeader)
	if err != nil {
		log.Errorf("websocketproxy: couldn't upgrade %s", err)
		return
	}
	defer connPub.Close()

	errBackend := make(chan error, 1)
	replicateWebsocketConn := func(dst, src *websocket.Conn, errc chan error) {
		for {
			msgType, msg, err := src.ReadMessage()
			if err != nil {
				m := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("%v", err))
				if e, ok := err.(*websocket.CloseError); ok {
					if e.Code != websocket.CloseNoStatusReceived {
						m = websocket.FormatCloseMessage(e.Code, e.Text)
					}
				}
				errc <- err

				dst.WriteMessage(websocket.CloseMessage, m)
				break
			}
			err = dst.WriteMessage(msgType, msg)
			if err != nil {
				errc <- err
				break
			}
		}
	}

	go func() {
		for c.status != Closed {
			msg, ok := <-c.queue
			if !ok {
				return
			}
			err = connPub.WriteMessage(msg.Type, msg.Value)
			if err != nil {
				break
			}
		}
	}()

	go replicateWebsocketConn(c.ws, connPub, errBackend)

	var message string
	select {
	case err = <-errBackend:
		message = "websocketproxy: Error when copying from client to backend: %v"
	}
	if e, ok := err.(*websocket.CloseError); !ok || e.Code == websocket.CloseAbnormalClosure {
		log.Errorf(message, err)
	}

	return nil
}

// Proxy a HTTP request through the Proxy over the websocket connection
func (c *Connection) proxyRequest(w http.ResponseWriter, r *http.Request) (err error) {
	log.Infof("proxy request to %s", c.pool.id)
	defer c.Release()

	// [1]: Serialize HTTP request
	jsonReq, err := json.Marshal(common.SerializeHTTPRequest(r))
	if err != nil {
		return fmt.Errorf("unable to serialize request : %w", err)
	}
	// i.e.
	// {
	// 		"Method":"GET",
	// 		"URL":"http://localhost:8081/hello",
	// 		"Header":{"Accept":["*/*"],"User-Agent":["curl/7.77.0"],"X-Proxy-Destination":["http://localhost:8081/hello"]},
	//		"ContentLength":0
	// }

	// [2]: Send the HTTP request to the peer
	// Send the serialized HTTP request to the the peer
	if err = c.ws.WriteMessage(websocket.BinaryMessage, jsonReq); err != nil {
		return fmt.Errorf("unable to write request : %w", err)
	}

	// Pipe the HTTP request body to the peer
	//bodyWriter, err := c.ws.NextWriter(websocket.BinaryMessage)
	//if err != nil {
	//	return fmt.Errorf("unable to get request body writer : %w", err)
	//}
	//if _, err = io.Copy(bodyWriter, r.Body); err != nil {
	//	return fmt.Errorf("unable to pipe request body : %w", err)
	//}
	//if err = bodyWriter.Close(); err != nil {
	//	return fmt.Errorf("unable to pipe request body (close) : %w", err)
	//}

	msg, ok := <-c.queue
	if !ok {
		return
	}

	jsonResponse := msg.Value

	// Deserialize the HTTP Response
	httpResponse := &common.HTTPResponse{}
	if err = json.Unmarshal(jsonResponse, httpResponse); err != nil {
		return fmt.Errorf("unable to unserialize http response : %w", err)
	}

	// Write response headers back to the client
	for header, values := range httpResponse.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(httpResponse.StatusCode)

	//msg, ok = <-c.queue
	//if !ok {
	//	return
	//}
	//
	//responseBody := msg.Value
	//

	responseBodyReader := bytes.NewReader(httpResponse.Body)
	if _, err = io.Copy(w, responseBodyReader); err != nil {
		return fmt.Errorf("unable to pipe response body : %w", err)
	}

	return
}

// Take notifies that this connection is going to be used
func (c *Connection) Take() bool {

	if c.status == Closed || c.status == Busy {
		return false
	}

	c.status = Busy
	return true
}

// Release notifies that this connection is ready to use again
func (c *Connection) Release() {

	if c.status == Closed {
		return
	}

	c.idleSince = time.Now()
	c.status = Idle

	go c.pool.Offer(c)
}

// Close the connection
func (c *Connection) Close() {
	if c.status == Closed {
		return
	}
	c.status = Closed

	log.Infof("Closing connection from %s", c.pool.id)

	// Close the underlying TCP connection
	c.ws.Close()
}
