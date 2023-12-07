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
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/e154/smart-home/system/gate/common"
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
	lock      sync.Mutex
}

// NewConnection returns a new Connection.
func NewConnection(pool *Pool, ws *websocket.Conn) *Connection {
	c := &Connection{
		pool:   pool,
		ws:     ws,
		status: Idle,
	}

	// Mark that this connection is ready to use for relay
	c.Release()

	return c
}

func (c *Connection) proxyWs(w http.ResponseWriter, r *http.Request) (err error) {

	// Only pass those headers to the upgrader.
	//upgradeHeader := http.Header{}
	//if hdr := r.Header.Get("Sec-Websocket-Protocol"); hdr != "" {
	//	upgradeHeader.Set("Sec-Websocket-Protocol", hdr)
	//}
	//if hdr := r.Header.Get("Set-Cookie"); hdr != "" {
	//	upgradeHeader.Set("Set-Cookie", hdr)
	//}

	_ = c.ws.WriteMessage(websocket.TextMessage, []byte("WS"))

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	// Now upgrade the existing incoming request to a WebSocket connection.
	// Also pass the header that we gathered from the Dial handshake.
	connPub, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("websocketproxy: couldn't upgrade %s", err)
		return
	}
	defer connPub.Close()

	errClient := make(chan error, 1)
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

	go replicateWebsocketConn(connPub, c.ws, errClient)
	go replicateWebsocketConn(c.ws, connPub, errBackend)

	var message string
	select {
	case err = <-errClient:
		message = "websocketproxy: Error when copying from backend to client: %v"
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
	if err := c.ws.WriteMessage(websocket.TextMessage, jsonReq); err != nil {
		return fmt.Errorf("unable to write request : %w", err)
	}

	// Pipe the HTTP request body to the peer
	bodyWriter, err := c.ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return fmt.Errorf("unable to get request body writer : %w", err)
	}
	if _, err := io.Copy(bodyWriter, r.Body); err != nil {
		return fmt.Errorf("unable to pipe request body : %w", err)
	}
	if err := bodyWriter.Close(); err != nil {
		return fmt.Errorf("unable to pipe request body (close) : %w", err)
	}

	// [3]: Wait the HTTP response is ready
	//responseChannel := make(chan (io.Reader))
	//c.nextResponse <- responseChannel
	//responseReader, ok := <-responseChannel
	//if responseReader == nil {
	//	if ok {
	//		// The value of ok is false, the channel is closed and empty.
	//		// See the Receiver operator in https://go.dev/ref/spec for more information.
	//		close(responseChannel)
	//	}
	//	return fmt.Errorf("unable to get http response reader : %w", err)
	//}
	//
	//// [4]: Read the HTTP response from the peer
	//// Get the serialized HTTP Response from the peer
	//jsonResponse, err := io.ReadAll(responseReader)
	//if err != nil {
	//	close(responseChannel)
	//	return fmt.Errorf("unable to read http response : %w", err)
	//}
	//
	//// Notify the read() goroutine that we are done reading the response
	//close(responseChannel)

	messageType, jsonResponse, err := c.ws.ReadMessage()
	if messageType == -1 {
		return
	}

	if err != nil {
		log.Error(err.Error())
		return
	}

	// Deserialize the HTTP Response
	httpResponse := new(common.HTTPResponse)
	if err := json.Unmarshal(jsonResponse, httpResponse); err != nil {
		return fmt.Errorf("unable to unserialize http response : %w", err)
	}

	// Write response headers back to the client
	for header, values := range httpResponse.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(httpResponse.StatusCode)

	// [5]: Wait the HTTP response body is ready
	// Get the HTTP Response body from the the peer
	// To do so send a new channel to the read() goroutine
	// to get the next message reader
	//responseBodyChannel := make(chan (io.Reader))
	//c.nextResponse <- responseBodyChannel
	//responseBodyReader, ok := <-responseBodyChannel
	//if responseBodyReader == nil {
	//	if ok {
	//		// If more is false the channel is already closed
	//		close(responseChannel)
	//	}
	//	return fmt.Errorf("unable to get http response body reader : %w", err)
	//}

	messageType, responseBody, err := c.ws.ReadMessage()
	if messageType == -1 {
		return
	}

	responseBodyReader := bytes.NewReader(responseBody)

	// [6]: Read the HTTP response body from the peer
	// Pipe the HTTP response body right from the remote Proxy to the client
	if _, err := io.Copy(w, responseBodyReader); err != nil {
		//close(responseBodyChannel)
		return fmt.Errorf("unable to pipe response body : %w", err)
	}

	// Notify read() that we are done reading the response body
	//close(responseBodyChannel)

	c.Release()

	return
}

// Take notifies that this connection is going to be used
func (c *Connection) Take() bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.status == Closed {
		return false
	}

	if c.status == Busy {
		return false
	}

	c.status = Busy
	return true
}

// Release notifies that this connection is ready to use again
func (c *Connection) Release() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.status == Closed {
		return
	}

	c.idleSince = time.Now()
	c.status = Idle

	go c.pool.Offer(c)
}

// Close the connection
func (c *Connection) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.close()
}

// Close the connection ( without lock )
func (c *Connection) close() {
	if c.status == Closed {
		return
	}

	log.Infof("Closing connection from %s", c.pool.id)

	// This one will be executed *before* lock.Unlock()
	defer func() { c.status = Closed }()

	// Unlock a possible read() wild message
	//close(c.nextResponse)

	// Close the underlying TCP connection
	c.ws.Close()
}
