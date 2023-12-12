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
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/gate/common"
)

var (
	log = logger.MustGetLogger("wsp")
)

// Server is a Reverse HTTP Proxy over WebSocket
// This is the Server part, Clients will offer websocket connections,
// those will be pooled to transfer HTTP Request and response
type Server struct {
	Config *Config

	upgrader websocket.Upgrader

	// In pools, keep connections with WebSocket peers.
	pools map[PoolID]*Pool

	// A RWMutex is a reader/writer mutual exclusion lock,
	// and it is for exclusive control with pools operation.
	//
	// This is locked when reading and writing pools, the timing is when:
	// 1. (rw) registering websocket clients in /register endpoint
	// 2. (rw) remove empty pools which has no connections
	// 3. (r) dispatching connection from available pools to clients requests
	//
	// And then it is released after each process is completed.
	lock sync.RWMutex
	done chan struct{}

	// Through dispatcher channel it communicates between "server" thread and "dispatcher" thread.
	// "server" thread sends the value to this channel when accepting requests in the endpoint /requests,
	// and "dispatcher" thread reads this channel.
	dispatcher chan *ConnectionRequest

	server *http.Server
}

// ConnectionRequest is used to request a proxy connection from the dispatcher
type ConnectionRequest struct {
	connection chan *Connection
	serverId   PoolID
}

// NewConnectionRequest creates a new connection request
func NewConnectionRequest(timeout time.Duration, serverId PoolID) (cr *ConnectionRequest) {
	cr = &ConnectionRequest{
		connection: make(chan *Connection),
		serverId:   serverId,
	}
	return
}

// NewServer return a new Server instance
func NewServer(config *Config) (server *Server) {
	server = &Server{
		Config:     config,
		upgrader:   websocket.Upgrader{},
		lock:       sync.RWMutex{},
		done:       make(chan struct{}),
		dispatcher: make(chan *ConnectionRequest),
		pools:      make(map[PoolID]*Pool),
	}
	return
}

// Start Server HTTP server
func (s *Server) Start() {
	go func() {
	L:
		for {
			select {
			case <-s.done:
				break L
			case <-time.After(5 * time.Second):
				s.clean()
			}
		}
	}()

	//r := http.NewServeMux()
	// TODO: I want to detach the handler function from the Server struct,
	// but it is tightly coupled to the internal state of the Server.
	//r.HandleFunc("/gate/register", s.Register)
	//r.HandleFunc("/gate/request", s.Request)
	//r.HandleFunc("/gate/status", s.status)

	// Dispatch connection from available pools to clients requests
	// in a separate thread from the server thread.
	go s.dispatchConnections()

	//s.server = &http.Server{
	//	Addr:    s.Config.GetAddr(),
	//	Handler: r,
	//}
	//go s.server.ListenAndServe()
}

// clean removes empty Pools which has no connection.
// It is invoked every 5 sesconds and at shutdown.
func (s *Server) clean() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.pools) == 0 {
		return
	}

	idle := 0
	busy := 0

	for _, pool := range s.pools {
		if pool.IsEmpty() {
			log.Infof("Removing empty connection pool : %s", pool.id)
			pool.Shutdown()
			delete(s.pools, pool.id)
		}

		ps := pool.Size()
		idle += ps.Idle
		busy += ps.Busy
	}

	log.Infof("%d pools, %d idle, %d busy", len(s.pools), idle, busy)
}

// Dispatch connection from available pools to clients requests
func (s *Server) dispatchConnections() {
	for {
		// Runs in an infinite loop and keeps receiving the value from the `server.dispatcher` channel
		// The operator <- is "receive operator", which expression blocks until a value is available.
		request, ok := <-s.dispatcher
		if !ok {
			// The value of `ok` is false if it is a zero value generated because the channel is closed an empty.
			// In this case, that means server shutdowns.
			break
		}

		// A timeout is set for each dispatch request.
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, s.Config.GetTimeout())
		defer cancel()

	L:
		for {
			select {
			case <-ctx.Done(): // The timeout elapses
				break L
			default: // Go through
			}

			s.lock.RLock()
			if len(s.pools) == 0 {
				// No connection pool available
				s.lock.RUnlock()
				break
			}

			if request.serverId != "" {
				pool, ok := s.pools[request.serverId]
				if !ok {
					log.Warnf("pool %s not found", request.serverId)
					s.lock.RUnlock()
					break
				}
				s.lock.RUnlock()
				connection, ok := <- pool.idle
				if !ok {
					continue // a pool has been removed, try again
				}
				// [2]: Verify that we can use this connection and take it.
				if connection.Take() {
					request.connection <- connection
					break
				}
				continue
			} else {
				if len(s.pools) > 1 {
					s.lock.RUnlock()
					break
				}
			}

			// [1]: Select a pool which has an idle connection
			// Build a select statement dynamically to handle an arbitrary number of pools.
			cases := make([]reflect.SelectCase, len(s.pools)+1)
			var i int
			for _, pool := range s.pools {
				cases[i] = reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(pool.idle),
				}
				i++
			}
			cases[len(cases)-1] = reflect.SelectCase{
				Dir: reflect.SelectDefault,
			}
			s.lock.RUnlock()

			_, value, ok := reflect.Select(cases)
			if !ok {
				continue // a pool has been removed, try again
			}
			connection, _ := value.Interface().(*Connection)

			// [2]: Verify that we can use this connection and take it.
			if connection.Take() {
				request.connection <- connection
				break
			}
		}

		close(request.connection)
	}
}

func (s *Server) Ws(w http.ResponseWriter, r *http.Request) {
	log.Infof("[%s] %s", r.Method, r.URL.String())

	if len(s.pools) == 0 {
		common.ProxyErrorf(w, "No proxy available")
		return
	}

	serverId, err := s.GetServerID(r)
	if err != nil {
		common.ProxyErrorf(w, err.Error())
		return
	}

	// [2]: Take an WebSocket connection available from pools for relaying received requests.
	request := NewConnectionRequest(s.Config.GetTimeout(), PoolID(serverId))
	// "Dispatcher" is running in a separate thread from the server by `go s.dispatchConnections()`.
	// It waits to receive requests to dispatch connection from available pools to clients requests.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L93
	//
	// Notify request from handler to dispatcher through Server.dispatcher channel.
	s.dispatcher <- request
	// Dispatcher tries to find an available connection pool,
	// and it returns the connection through Server.connection channel.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L189
	//
	// Here waiting for a result from dispatcher.
	connection := <-request.connection
	if connection == nil {
		// It means that dispatcher has set `nil` which is a system error case that is
		// not expected in the normal flow.
		common.ProxyErrorf(w, "Unable to get a proxy connection")
		return
	}
	defer func() {
		connection.Close()
	}()

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
	// [1]: Receive requests to be proxied
	// Parse destination URL
	//dstURL := r.Header.Get("X-PROXY-DESTINATION")
	//if dstURL == "" {
	//	common.ProxyErrorf(w, "Missing X-PROXY-DESTINATION header")
	//	return
	//}
	//URL, err := url.Parse(dstURL)
	//if err != nil {
	//	common.ProxyErrorf(w, "Unable to parse X-PROXY-DESTINATION header")
	//	return
	//}
	//r.URL = URL

	if len(s.pools) == 0 {
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

	// [2]: Take an WebSocket connection available from pools for relaying received requests.
	request := NewConnectionRequest(s.Config.GetTimeout(), PoolID(serverId))
	// "Dispatcher" is running in a separate thread from the server by `go s.dispatchConnections()`.
	// It waits to receive requests to dispatch connection from available pools to clients requests.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L93
	//
	// Notify request from handler to dispatcher through Server.dispatcher channel.
	s.dispatcher <- request
	// Dispatcher tries to find an available connection pool,
	// and it returns the connection through Server.connection channel.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L189
	//
	// Here waiting for a result from dispatcher.
	connection := <-request.connection
	if connection == nil {
		// It means that dispatcher has set `nil` which is a system error case that is
		// not expected in the normal flow.
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
	if s.Config.SecretKey != "" && secretKey != s.Config.SecretKey {
		common.ProxyErrorf(w, "Invalid X-SECRET-KEY")
		return
	}

	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		common.ProxyErrorf(w, "HTTP upgrade error : %v", err)
		return
	}

	// 2. Wait a greeting message from the peer and parse it
	// The first message should contains the remote Proxy name and size
	_, greeting, err := ws.ReadMessage()
	if err != nil {
		common.ProxyErrorf(w, "Unable to read greeting message : %s", err)
		ws.Close()
		return
	}

	// Parse the greeting message
	split := strings.Split(string(greeting), "_")
	id := PoolID(split[0])
	size, err := strconv.Atoi(split[1])
	if err != nil {
		common.ProxyErrorf(w, "Unable to parse greeting message : %s", err)
		ws.Close()
		return
	}

	// 3. Register the connection into server pools.
	// s.lock is for exclusive control of pools operation.
	s.lock.Lock()
	defer s.lock.Unlock()

	var pool *Pool
	var ok bool
	if pool, ok = s.pools[id]; !ok {
		pool = NewPool(s, id)
		s.pools[id] = pool
	}

	// update pool size
	pool.size = size

	// Add the WebSocket connection to the pool
	pool.Register(ws)
}

// Shutdown stop the Server
func (s *Server) Shutdown() {
	close(s.done)
	close(s.dispatcher)
	for _, pool := range s.pools {
		pool.Shutdown()
	}
	s.clean()
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
