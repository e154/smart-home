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
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/system/stream"
	"sync"
	"time"
)

// Pool manage a pool of connection to a remote Server
type Pool struct {
	client    *Client
	target    string
	secretKey string

	connections []*Connection
	lock        sync.RWMutex

	done chan struct{}

	api    *api.Api
	stream *stream.Stream
	status Status
}

// NewPool creates a new Pool
func NewPool(client *Client, target string, secretKey string, api *api.Api, stream *stream.Stream) *Pool {
	return &Pool{
		client:      client,
		target:      target,
		connections: make([]*Connection, 0),
		secretKey:   secretKey,
		done:        make(chan struct{}),
		api:         api,
		stream:      stream,
	}
}

// Start connect to the remote Server
func (p *Pool) Start(ctx context.Context) {
	p.connector(ctx)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

	L:
		for {
			select {
			case <-p.done:
				break L
			case <-ticker.C:
				p.connector(ctx)
			}
		}
	}()
}

// The garbage collector
func (p *Pool) connector(ctx context.Context) {
	p.lock.Lock()
	defer p.lock.Unlock()

	poolSize := p.Size()
	p.status = Status{
		Connecting: poolSize.connecting,
		Idle:       poolSize.idle,
		Running:    poolSize.running,
		Total:      poolSize.total,
	}

	// Create enough connection to fill the pool
	toCreate := p.client.cfg.PoolIdleSize - poolSize.idle

	// Create only one connection if the pool is empty
	if poolSize.total == 0 {
		toCreate = 1
	}

	// Ensure to open at most PoolMaxSize connections
	if poolSize.total+toCreate > p.client.cfg.PoolMaxSize {
		toCreate = p.client.cfg.PoolMaxSize - poolSize.total
	}

	// Try to reach ideal p size
	for i := 0; i < toCreate; i++ {
		conn := NewConnection(p, p.api, p.stream)
		p.connections = append(p.connections, conn)

		go func() {
			err := conn.Connect(ctx)
			if err != nil {
				//log.Errorf("Unable to connect to %s : %s", p.target, err)

				p.lock.Lock()
				defer p.lock.Unlock()
				p.remove(conn)
			}
		}()
	}
}

// Add a connection to the pool
func (p *Pool) add(conn *Connection) {
	p.connections = append(p.connections, conn)
}

// Remove a connection from the pool
func (p *Pool) remove(conn *Connection) {
	// This trick uses the fact that a slice shares the same backing array and capacity as the original,
	// so the storage is reused for the filtered slice. Of course, the original contents are modified.

	var filtered []*Connection // == nil
	for _, c := range p.connections {
		if conn != c {
			filtered = append(filtered, c)
		}
	}
	p.connections = filtered
}

// Shutdown close all connection in the pool
func (p *Pool) Shutdown() {
	close(p.done)
	for _, conn := range p.connections {
		conn.Close()
	}
}

// Size return the current state of the pool
func (p *Pool) Size() (poolSize *PoolSize) {
	poolSize = &PoolSize{}
	poolSize.total = len(p.connections)
	for _, connection := range p.connections {
		switch connection.status {
		case CONNECTING:
			poolSize.connecting++
		case IDLE:
			poolSize.idle++
		case RUNNING:
			poolSize.running++
		}
	}

	return
}
