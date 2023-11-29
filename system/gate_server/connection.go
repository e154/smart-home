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

package gate_server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/e154/smart-home/system/gate_client/common"
)

// ConnectionStatus is an enumeration type which represents the status of WebSocket connection.
type ConnectionStatus int

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
	// nextResponse is the channel of channel to wait an HTTP response.
	//
	// In advance, the `read` function waits to receive the HTTP response as a separate thread "reader".
	// (See https://github.com/hgsgtk/wsp/blob/29cc73bbd67de18f1df295809166a7a5ef52e9fa/server/connection.go#L56 )
	//
	// When a "server" thread proxies, it sends the HTTP request to the peer over the WebSocket,
	// and sends the channel of the io.Reader interface (chan io.Reader) that can read the HTTP response to the field `nextResponse`,
	// then waits until the value is written in the channel (chan io.Reader) by another thread "reader".
	//
	// After the thread "reader" detects that the HTTP response from the peer of the WebSocket connection has been written,
	// it sends the value to the channel (chan io.Reader),
	// and the "server" thread can proceed to process the rest procedures.
	nextResponse chan chan io.Reader
}

// NewConnection returns a new Connection.
func NewConnection(pool *Pool, ws *websocket.Conn) *Connection {
	// Initialize a new Connection
	c := new(Connection)
	c.pool = pool
	c.ws = ws
	c.nextResponse = make(chan chan io.Reader)
	c.status = Idle

	// Mark that this connection is ready to use for relay
	c.Release()

	// Start to listen to incoming messages over the WebSocket connection
	go c.read()

	return c
}

// read the incoming message of the connection
func (c *Connection) read() {
	defer func() {
		if r := recover(); r != nil {
			log.Infof("Websocket crash recovered : %s", r)
		}
		c.Close()
	}()

	for {
		if c.status == Closed {
			break
		}

		// https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
		//
		// We need to ensure :
		//  - no concurrent calls to ws.NextReader() / ws.ReadMessage()
		//  - only one reader exists at a time
		//  - wait for reader to be consumed before requesting the next one
		//  - always be reading on the socket to be able to process control messages ( ping / pong / close )

		// We will block here until a message is received or the ws is closed
		_, reader, err := c.ws.NextReader()
		if err != nil {
			break
		}

		if c.status != Busy {
			// We received a wild unexpected message
			break
		}

		// When it gets here, it is expected to be either a HttpResponse or a HttpResponseBody has been returned.
		//
		// Next, it waits to receive the value from the Connection.proxyRequest function that is invoked in the "server" thread.
		// https://github.com/hgsgtk/wsp/blob/29cc73bbd67de18f1df295809166a7a5ef52e9fa/server/connection.go#L157
		ch := <-c.nextResponse
		if ch == nil {
			// We have been unlocked by Close()
			break
		}

		// Send the reader back to Connection.proxyRequest
		ch <- reader

		// Wait for proxyRequest to close the channel
		// this notify that it is done with the reader
		<-ch
	}
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

	// Pipe the HTTP request body to the the peer
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
	responseChannel := make(chan (io.Reader))
	c.nextResponse <- responseChannel
	responseReader, ok := <-responseChannel
	if responseReader == nil {
		if ok {
			// The value of ok is false, the channel is closed and empty.
			// See the Receiver operator in https://go.dev/ref/spec for more information.
			close(responseChannel)
		}
		return fmt.Errorf("unable to get http response reader : %w", err)
	}

	// [4]: Read the HTTP response from the peer
	// Get the serialized HTTP Response from the peer
	jsonResponse, err := io.ReadAll(responseReader)
	if err != nil {
		close(responseChannel)
		return fmt.Errorf("unable to read http response : %w", err)
	}

	// Notify the read() goroutine that we are done reading the response
	close(responseChannel)

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
	responseBodyChannel := make(chan (io.Reader))
	c.nextResponse <- responseBodyChannel
	responseBodyReader, ok := <-responseBodyChannel
	if responseBodyReader == nil {
		if ok {
			// If more is false the channel is already closed
			close(responseChannel)
		}
		return fmt.Errorf("unable to get http response body reader : %w", err)
	}

	// [6]: Read the HTTP response body from the peer
	// Pipe the HTTP response body right from the remote Proxy to the client
	if _, err := io.Copy(w, responseBodyReader); err != nil {
		close(responseBodyChannel)
		return fmt.Errorf("unable to pipe response body : %w", err)
	}

	// Notify read() that we are done reading the response body
	close(responseBodyChannel)

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
	close(c.nextResponse)

	// Close the underlying TCP connection
	c.ws.Close()
}
