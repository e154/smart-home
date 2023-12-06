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
	"github.com/e154/smart-home/system/stream"
	"net/http"

	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/common/logger"
	"github.com/gorilla/websocket"
)

var (
	log = logger.MustGetLogger("wsp")
)

// Client connects to one or more Server using HTTP websockets.
// The Server can then send HTTP requests to execute.
type Client struct {
	cfg *Config

	client *http.Client
	dialer *websocket.Dialer
	pools  map[string]*Pool
	api    *api.Api
	stream *stream.Stream
}

// NewClient creates a new Client.
func NewClient(config *Config, api *api.Api, stream *stream.Stream) (c *Client) {
	c = &Client{
		cfg:    config,
		client: &http.Client{},
		dialer: &websocket.Dialer{},
		pools:  make(map[string]*Pool),
		api:    api,
		stream: stream,
	}
	return
}

// Start the Proxy
func (c *Client) Start(ctx context.Context) {
	for _, target := range c.cfg.Targets {
		pool := NewPool(c, target, c.cfg.SecretKey, c.api, c.stream)
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
