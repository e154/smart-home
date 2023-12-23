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
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/atomic"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = logger.MustGetLogger("wsp")
)

// Client connects to one or more Server using HTTP websockets.
// The Server can then send HTTP requests to execute.
type Client struct {
	cfg        *Config
	client     *http.Client
	dialer     *websocket.Dialer
	pools      map[string]*Pool
	api        *api.Api
	stream     *stream.Stream
	isStarted  *atomic.Bool
	adaptors   *adaptors.Adaptors
	jwtManager jwt_manager.JwtManager
}

// NewClient creates a new Client.
func NewClient(cfg *Config, api *api.Api, stream *stream.Stream, adaptors *adaptors.Adaptors, jwtManager jwt_manager.JwtManager) *Client {
	return &Client{
		cfg:        cfg,
		client:     &http.Client{},
		dialer:     &websocket.Dialer{},
		pools:      make(map[string]*Pool),
		api:        api,
		stream:     stream,
		isStarted:  atomic.NewBool(false),
		adaptors:   adaptors,
		jwtManager: jwtManager,
	}
}

// Start the Proxy
func (c *Client) Start(ctx context.Context) {
	//log.Info("Start")
	if !c.isStarted.CompareAndSwap(false, true) {
		return
	}
	for _, target := range c.cfg.Targets {
		pool := NewPool(c, target, c.cfg.SecretKey, c.api, c.stream, c.adaptors, c.jwtManager)
		c.pools[target] = pool
		go pool.Start(ctx)
	}
}

// Shutdown the Proxy
func (c *Client) Shutdown() {
	//log.Info("Shutdown")
	if !c.isStarted.CompareAndSwap(true, false) {
		return
	}
	for _, pool := range c.pools {
		pool.Shutdown()
	}
}
