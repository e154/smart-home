// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

// Client ...
type Client struct {
	Ip        string
	Referer   string
	UserAgent string
	Width     int
	Height    int
	Cookie    bool
	Language  string
	Platform  string
	Location  string
	Href      string
	Send      chan []byte // message buffered channel
	writeLock sync.Mutex
	Connect   *websocket.Conn
}

// UpdateInfo ...
func (c *Client) UpdateInfo(msg Message) {

	v := msg.Payload

	width, ok := v["width"].(float64)
	if ok {
		c.Width = int(width)
	}

	if height, ok := v["height"].(float64); ok {
		c.Height = int(height)
	}

	if cookie, ok := v["cookie"].(bool); ok {
		c.Cookie = cookie
	}

	if language, ok := v["language"].(string); ok {
		c.Language = language
	}

	if platform, ok := v["platform"].(string); ok {
		c.Platform = platform
	}

	if location, ok := v["location"].(string); ok {
		c.Location = location
	}

	if href, ok := v["href"].(string); ok {
		c.Href = href
	}
}

// Notify ...
func (c *Client) Notify(t, b string) {

	msg := &Message{
		Type:    Notify,
		Forward: Request,
		Payload: map[string]interface{}{
			"type": t,
			"body": b,
		},
	}

	c.Send <- msg.Pack()
}

func (c *Client) selfWrite(opCode int, payload []byte) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	_ = c.Connect.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Connect.WriteMessage(opCode, payload)
}

// send message to client
//
func (c *Client) WritePump() {

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		if c.Connect != nil {
			_ = c.Connect.Close()
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:

			if !ok {
				c.selfWrite(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.selfWrite(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.selfWrite(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Close ...
func (c *Client) Close() {

	err := c.selfWrite(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return
	}
}

// Write ...
func (c *Client) Write(payload []byte) (err error) {
	c.Send <- payload
	return
}
