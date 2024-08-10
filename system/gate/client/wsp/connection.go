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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/gate/common"
	"github.com/e154/smart-home/system/stream"
)

// Status of a Connection
const (
	CONNECTING = iota
	IDLE
	RUNNING
)

// Connection handle a single websocket (HTTP/TCP) connection to an Server
type Connection struct {
	pool   *Pool
	status int
	api    *api.Api
	stream *stream.Stream
	cli    *stream.Client
	debug  bool
	*sync.Mutex
	ws       *websocket.Conn
	isClosed *atomic.Bool
}

// NewConnection create a Connection object
func NewConnection(pool *Pool,
	api *api.Api,
	stream *stream.Stream) *Connection {
	c := &Connection{
		pool:     pool,
		status:   CONNECTING,
		api:      api,
		stream:   stream,
		Mutex:    &sync.Mutex{},
		isClosed: atomic.NewBool(true),
	}
	return c
}

// Connect to the IsolatorServer using a HTTP websocket
func (c *Connection) Connect(ctx context.Context) (err error) {
	//log.Infof("Connecting to %s", c.pool.target)

	// Create a new TCP(/TLS) connection ( no use of net.http )
	c.ws, _, err = c.pool.client.dialer.DialContext(
		ctx,
		c.pool.target,
		http.Header{"X-SECRET-KEY": {c.pool.secretKey}},
	)

	if err != nil {
		return err
	}

	c.ws.SetCloseHandler(func(code int, text string) error {
		c.isClosed.Store(true)
		message := websocket.FormatCloseMessage(code, "")
		return c.ws.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
	})

	c.isClosed.Store(false)

	if c.debug {
		log.Infof("Connected to %s", c.pool.target)
		defer log.Info("Connection closed ...")
	}

	// Send the greeting message with proxy id and wanted pool size.
	greeting := fmt.Sprintf(
		"%s_%d",
		c.pool.client.cfg.ID,
		c.pool.client.cfg.PoolIdleSize,
	)
	if err = c.WriteMessage(websocket.TextMessage, []byte(greeting)); err != nil {
		log.Error("greeting error :", err.Error())
		c.Close()
		return err
	}

	c.serve(ctx)

	return
}

// the main loop it :
//   - wait to receive HTTP requests from the Server
//   - execute HTTP requests
//   - send HTTP response back to the Server
//
// As in the server code there is no buffering of HTTP request/response body
// As is the server if any error occurs the connection is closed/throwed
func (c *Connection) serve(ctx context.Context) {

	// Keep connection alive
	go func() {
		timer := time.NewTicker(time.Second * 10)
		defer timer.Stop()
		for {
			select {
			case t := <-timer.C:
				err := c.ws.WriteControl(websocket.PingMessage, []byte{}, t.Add(time.Second))
				if err != nil {
					c.Close()
					return
				}
			}
		}
	}()

	var jsonRequest []byte
	var err error
	var messageType int

	for {
		// Read request
		c.status = IDLE
		messageType, jsonRequest, err = c.ws.ReadMessage()

		if messageType == -1 {
			break
		}

		if err != nil {
			log.Errorf("Unable to read request", err)
			break
		}

		c.status = RUNNING

		// Trigger a pool refresh to open new connections if needed
		go c.pool.connector(ctx)

		if strings.Contains(string(jsonRequest), "WS:") {
			// get user
			accessToken := string(jsonRequest)
			accessToken = strings.ReplaceAll(accessToken, "WS:", "")
			if accessToken == "" {
				log.Error(apperr.ErrUnauthorized.Error())
				return
			}
			var user *m.User
			user, err = c.pool.GetUser(accessToken)
			if err != nil {
				log.Error(apperr.ErrAccessDenied.Error())
				return
			}
			c.stream.NewConnection(c.ws, user)
			return
		}

		// Deserialize request
		req, err := common.DeserializeHTTPRequest(jsonRequest)
		if err != nil {
			c.error(fmt.Sprintf("Unable to deserialize http request : %v\n", err))
			break
		}

		//log.Infof("[%s] %s", req.Method, req.URL.String())

		req.RequestURI = req.URL.String()
		resp := httptest.NewRecorder()
		c.api.Echo().ServeHTTP(resp, req)

		// Serialize response
		jsonResponse, err := json.Marshal(common.SerializeHTTPResponse(resp))
		if err != nil {
			err = c.error(fmt.Sprintf("Unable to serialize response : %v\n", err))
			if err != nil {
				break
			}
			continue
		}

		// Write response
		err = c.WriteMessage(websocket.BinaryMessage, jsonResponse)
		if err != nil {
			log.Errorf("Unable to write response : %v", err)
			break
		}

	}
}

func (c *Connection) error(msg string) (err error) {
	resp := common.NewHTTPResponse()
	resp.StatusCode = 527

	log.Error(msg)

	resp.ContentLength = int64(len(msg))

	// Serialize response
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Unable to serialize response : %v", err)
		return
	}

	// Write response
	err = c.WriteMessage(websocket.TextMessage, jsonResponse)
	if err != nil {
		log.Errorf("Unable to write response : %v", err)
		return
	}

	// Write response body
	err = c.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		log.Errorf("Unable to write response body : %v", err)
		return
	}

	return
}

// Close close the ws/tcp connection and remove it from the pool
func (c *Connection) Close() {
	if !c.isClosed.CompareAndSwap(false, true) {
		return
	}

	c.isClosed.Store(true)

	if c.ws != nil {
		if err := c.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			//debug.PrintStack()
			//log.Error(err.Error())
		}
		c.ws.Close()
	}
}

func (c *Connection) WriteMessage(messageType int, data []byte) error {
	if c.isClosed.Load() {
		return errors.New("connection is closed")
	}

	c.Lock()
	defer func() {
		c.Unlock()
		if r := recover(); r != nil {
			log.Warn("Recovered")
			debug.PrintStack()
		}
	}()

	return c.ws.WriteMessage(messageType, data)
}
