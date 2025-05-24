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
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/e154/smart-home/internal/system/gate/common"
	"github.com/e154/smart-home/pkg/logger"

	"github.com/gorilla/websocket"
)

var (
	log = logger.MustGetLogger("wsp server")
)

// Server is a Reverse HTTP Proxy over WebSocket
// This is the Server part, Clients will offer websocket connections,
// those will be pooled to transfer HTTP Request and response
type Server struct {
	config   *Config
	upgrader websocket.Upgrader
	pools    *Pools
	done     chan struct{}
	server   *http.Server
}

func NewServer(config *Config) (server *Server) {
	server = &Server{
		config:   config,
		upgrader: websocket.Upgrader{},
		done:     make(chan struct{}),
		pools:    NewPools(config.Timeout, config.IdleTimeout),
	}
	return
}

// Start Server HTTP server
func (s *Server) Start() {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case <-time.After(5 * time.Second):
				s.pools.Clean()
			}
		}
	}()

}

func (s *Server) Ws(w http.ResponseWriter, r *http.Request) {
	log.Infof("[%s] %s", r.Method, r.URL.String())

	if s.pools.IsEmpty() {
		common.ProxyErrorf(w, "No proxy available")
		return
	}

	serverId, err := s.GetServerID(r)
	if err != nil {
		common.ProxyErrorf(w, err.Error())
		return
	}

	pool, ok := s.pools.GetPool(PoolID(serverId))
	if !ok {
		common.ProxyErrorf(w, "Unable to get a pool")
		return
	}

	connection := pool.GetIdleConnection(r.Context())
	if connection == nil {
		common.ProxyErrorf(w, "Unable to get a proxy connection")
		return
	}

	if err := connection.proxyWs(w, r); err != nil {
		// An error occurred throw the connection away
		log.Error(err.Error())
		connection.Close()

		// Try to return an error to the client
		// This might fail if response headers have already been sent
		common.ProxyError(w, err)
	}
}

func (s *Server) Request(w http.ResponseWriter, r *http.Request) {

	if s.pools.IsEmpty() {
		common.ProxyErrorf(w, "No proxy available")
		return
	}

	serverId, err := s.GetServerID(r)
	if err != nil {
		common.ProxyErrorf(w, err.Error())
		return
	}

	r.URL = &url.URL{
		Path:        r.URL.Path,
		RawQuery:    r.URL.RawQuery,
		Fragment:    r.URL.Fragment,
		RawFragment: r.URL.RawFragment,
	}

	log.Infof("[%s] %s", r.Method, r.URL.String())

	pool, ok := s.pools.GetPool(PoolID(serverId))
	if !ok {
		common.ProxyErrorf(w, "Unable to get a proxy connection")
		return
	}

	connection := pool.GetIdleConnection(r.Context())
	if connection == nil {
		common.ProxyErrorf(w, "Unable to get a proxy connection")
		return
	}

	// [3]: Send the request to the peer through the WebSocket connection.
	if err := connection.proxyRequest(w, r); err != nil {
		// An error occurred throw the connection away
		log.Error(err.Error())
		connection.Close()

		// Try to return an error to the client
		// This might fail if response headers have already been sent
		common.ProxyError(w, err)
	}
}

// Request receives the WebSocket upgrade handshake request from wsp_client.
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade a received HTTP request to a WebSocket connection
	secretKey := r.Header.Get("X-SECRET-KEY")
	if s.config.SecretKey != "" && secretKey != s.config.SecretKey {
		common.ProxyErrorf(w, "Invalid X-SECRET-KEY")
		return
	}

	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.ProxyErrorf(w, "HTTP upgrade error : %v", err)
		return
	}

	if err = s.pools.RegisterConnection(ws); err != nil {
		common.ProxyErrorf(w, err.Error())
	}
}

// Shutdown stop the Server
func (s *Server) Shutdown() {
	close(s.done)
	s.pools.Shutdown()
}

func (s *Server) GetServerID(r *http.Request) (serverID string, err error) {
	serverID = r.Header.Get("X-SERVER-ID")
	if serverID != "" {
		return
	}

	query := r.URL.Query()
	serverID = query.Get("serverId")
	if serverID != "" {
		return
	}

	serverID = query.Get("server_id")
	if serverID != "" {
		return
	}

	if serverID == "" {
		err = errors.New("Unable to parse DESTINATION params")
	}

	return
}
