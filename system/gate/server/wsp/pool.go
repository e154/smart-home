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
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Pool handles all connections from the peer.
type Pool struct {
	server *Server
	id     PoolID

	size int
	idle chan *Connection

	done        bool
	lock        sync.RWMutex
	connections map[string]*Connection
}

// PoolID represents the identifier of the connected WebSocket client.
type PoolID string

// NewPool creates a new Pool
func NewPool(server *Server, id PoolID) *Pool {
	p := &Pool{
		server:      server,
		id:          id,
		idle:        make(chan *Connection),
		connections: make(map[string]*Connection),
	}
	return p
}

// Register creates a new Connection and adds it to the pool
func (p *Pool) Register(ws *websocket.Conn) {
	p.lock.Lock()
	defer p.lock.Unlock()

	// Ensure we never add a connection to a pool we have garbage collected
	if p.done {
		return
	}

	log.Infof("Registering new connection from %s", p.id)
	connection := NewConnection(p, ws)
	id := uuid.NewString()

	go func() {
		p.connections[id] = connection
		connection.WritePump()
		delete(p.connections, id)
	}()

}

// Offer offers an idle connection to the server.
func (p *Pool) Offer(connection *Connection) {
	// The original code of root-gg/wsp was invoking goroutine,
	// but the callder was also invoking goroutine,
	// so it was deemed unnecessary and removed.
	p.idle <- connection
}

// Clean removes dead connection from the pool
// Look for dead connection in the pool
// This MUST be surrounded by pool.lock.Lock()
func (p *Pool) Clean() {
	idle := 0
	now := time.Now()
	for id, connection := range p.connections {
		// We need to be sur we'll never close a BUSY or soon to be BUSY connection
		connection.lock.Lock()
		if connection.status == Idle {
			idle++
			if idle > p.size+1 {
				// We have enough idle connections in the pool.
				// Terminate the connection if it is idle since more that IdleTimeout
				if now.Sub(connection.idleSince).Seconds() > p.server.Config.IdleTimeout.Seconds() {
					connection.close()
					delete(p.connections, id)
				}
			}
		}
		connection.lock.Unlock()
		if connection.status == Closed {
			connection.Close()
		}
	}
}

// IsEmpty clean the pool and return true if the pool is empty
func (p *Pool) IsEmpty() bool {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.Clean()
	return len(p.connections) == 0
}

// Shutdown closes every connections in the pool and cleans it
func (p *Pool) Shutdown() {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.done = true

	for id, connection := range p.connections {
		connection.Close()
		delete(p.connections, id)
	}
}

// Size return the number of connection in each state in the pool
func (p *Pool) Size() (ps *PoolSize) {
	p.lock.Lock()
	defer p.lock.Unlock()

	ps = &PoolSize{}
	for _, connection := range p.connections {
		if connection.status == Idle {
			ps.Idle++
		} else if connection.status == Busy {
			ps.Busy++
		} else if connection.status == Closed {
			ps.Closed++
		}
	}

	return
}
