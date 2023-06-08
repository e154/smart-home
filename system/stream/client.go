// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"io"

	"github.com/e154/smart-home/api/stub/api"
	m "github.com/e154/smart-home/models"
)

// Client ...
type Client struct {
	closed bool
	server api.StreamService_SubscribeServer
	user   *m.User
}

// NewClient ...
func NewClient(server api.StreamService_SubscribeServer, user *m.User) *Client {
	return &Client{server: server, user: user}
}

// WritePump ...
func (c *Client) WritePump(f func(*Client, string, string, []byte)) (err error) {

	var in *api.Request
	for !c.closed {
		in, err = c.server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		f(c, in.Id, in.Query, in.Body)
	}

	return
}

// Close ...
func (c *Client) Close() {
	c.closed = true
}

// Send ...
func (c *Client) Send(id, query string, body []byte) (err error) {
	err = c.server.Send(&api.Response{
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
