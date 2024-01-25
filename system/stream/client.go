// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package stream

import (
	"encoding/json"
	"sync"

	m "github.com/e154/smart-home/models"
	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	closed    bool
	user      *m.User
	sessionId string
	*sync.Mutex
	ws *websocket.Conn
}

// NewClient ...
func NewClient(ws *websocket.Conn, user *m.User, sessionId string) *Client {
	return &Client{
		ws:        ws,
		user:      user,
		Mutex:     &sync.Mutex{},
		sessionId: sessionId,
	}
}

// WritePump ...
func (c *Client) WritePump(f func(*Client, string, string, []byte)) {

	var data []byte
	var messageType int
	var err error
	for {
		messageType, data, err = c.ws.ReadMessage()
		if messageType == -1 || err != nil {
			return
		}

		msg := &Message{}
		_ = json.Unmarshal(data, msg)
		f(c, msg.Id, msg.Query, msg.Body)
	}
}

// Close ...
func (c *Client) Close() {
	c.closed = true
	if c.ws != nil {
		_ = c.ws.WriteMessage(websocket.CloseMessage, []byte{})
		_ = c.ws.Close()
	}
}

// Send ...
func (c *Client) Send(id, query string, body []byte) (err error) {
	c.Lock()
	defer c.Unlock()

	if c.closed {
		return
	}

	err = c.ws.WriteJSON(&Message{
		Id:    id,
		Query: query,
		Body:  body,
	})
	if err != nil {
		c.closed = true
	}
	return
}

func (c *Client) GetUser() *m.User {
	return c.user
}

func (c *Client) SessionID() string {
	return c.sessionId
}
