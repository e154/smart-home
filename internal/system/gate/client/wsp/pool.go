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
	"sync"
	"time"

	"github.com/e154/smart-home/internal/api"
	"github.com/e154/smart-home/internal/system/jwt_manager"
	"github.com/e154/smart-home/internal/system/stream"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"

	"github.com/google/uuid"
)

// Pool manage a pool of connection to a remote Server
type Pool struct {
	client    *Client
	target    string
	secretKey string

	connections sync.Map
	lock        sync.RWMutex

	done chan struct{}

	api    *api.Api
	stream *stream.Stream
	status Status

	adaptors   *adaptors.Adaptors
	jwtManager jwt_manager.JwtManager
}

// NewPool creates a new Pool
func NewPool(client *Client, target string,
	secretKey string,
	api *api.Api,
	stream *stream.Stream,
	adaptors *adaptors.Adaptors,
	jwtManager jwt_manager.JwtManager) *Pool {
	return &Pool{
		client:      client,
		target:      target,
		connections: sync.Map{},
		secretKey:   secretKey,
		done:        make(chan struct{}),
		api:         api,
		stream:      stream,
		adaptors:    adaptors,
		jwtManager:  jwtManager,
	}
}

// Start connect to the remote Server
func (p *Pool) Start(ctx context.Context) {
	//log.Info("Start")
	p.connector(ctx)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-p.done:
				return
			case <-ticker.C:
				p.connector(ctx)
			}
		}
	}()
}

// Shutdown close all connection in the pool
func (p *Pool) Shutdown() {
	//log.Info("Shutdown")
	close(p.done)
	p.connections.Range(func(key, value interface{}) bool {
		connection := value.(*Connection)
		connection.Close()
		p.connections.Delete(key)
		return true
	})
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

	if toCreate < 0 {
		toCreate = 0
	}

	// Create only one connection if the pool is empty
	if poolSize.total == 0 {
		toCreate = 1
	}

	// Ensure to open at most PoolMaxSize connections
	if poolSize.total+toCreate >= p.client.cfg.PoolMaxSize {
		toCreate = 0
	}

	if toCreate < 0 {
		toCreate = 0
	}

	if toCreate == 0 {
		return
	}

	// Try to reach ideal p size
	for i := 0; i < toCreate; i++ {
		connection := NewConnection(p, p.api, p.stream)
		id := uuid.NewString()
		p.connections.Store(id, connection)

		go func() {
			err := connection.Connect(ctx)
			if err != nil {
				//log.Errorf("Unable to connect to %s : %s", p.target, err)
			}
			p.connections.Delete(id)
		}()
	}
}

// Size return the current state of the pool
func (p *Pool) Size() (poolSize *PoolSize) {
	poolSize = &PoolSize{}
	p.connections.Range(func(key, value interface{}) bool {
		poolSize.total++
		connection := value.(*Connection)
		switch connection.status {
		case CONNECTING:
			poolSize.connecting++
		case IDLE:
			poolSize.idle++
		case RUNNING:
			poolSize.running++
		}
		return true
	})
	return
}

func (p *Pool) GetUser(accessToken string) (user *m.User, err error) {

	claims, err := p.jwtManager.Verify(accessToken)
	if err != nil {
		return
	}

	user, err = p.adaptors.User.GetById(context.Background(), claims.UserId)
	if err != nil {
		return
	}

	return
}
