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

package gate_client

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client connects to one or more Server using HTTP websockets.
// The Server can then send HTTP requests to execute.
type Client struct {
	Config *Config

	client *http.Client
	dialer *websocket.Dialer
	pools  map[string]*Pool
}

// NewClient creates a new Client.
func NewClient(config *Config) (c *Client) {
	c = new(Client)
	c.Config = config
	c.client = &http.Client{}
	c.dialer = &websocket.Dialer{}
	c.pools = make(map[string]*Pool)
	return
}

// Start the Proxy
func (c *Client) Start(ctx context.Context) {
	for _, target := range c.Config.Targets {
		pool := NewPool(c, target, c.Config.SecretKey)
		c.pools[target] = pool
		go pool.Start(ctx)
	}
}

// Shutdown the Proxy
func (c *Client) Shutdown() {
	for _, pool := range c.pools {
		pool.Shutdown()
	}
}
