// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Pools struct {
	timeout     time.Duration
	idleTimeout time.Duration
	*sync.Mutex
	pools map[PoolID]*Pool
}

func NewPools(timeout, idleTimeout time.Duration) *Pools {
	return &Pools{
		timeout:     timeout,
		idleTimeout: idleTimeout,
		Mutex:       &sync.Mutex{},
		pools:       make(map[PoolID]*Pool),
	}
}

func (p *Pools) Shutdown() {
	for _, pool := range p.pools {
		pool.Shutdown()
	}
}

func (p *Pools) RegisterConnection(ws *websocket.Conn) (err error) {

	// 2. Wait a greeting message from the peer and parse it
	// The first message should contains the remote Proxy name and size
	_, greeting, err := ws.ReadMessage()
	if err != nil {
		err = fmt.Errorf("Unable to read greeting message : %s", err)
		return
	}

	// Parse the greeting message
	split := strings.Split(string(greeting), "_")
	poolID := PoolID(split[0])
	size, err := strconv.Atoi(split[1])
	if err != nil {
		err = fmt.Errorf("Unable to parse greeting message : %s", err)
		return
	}

	p.Lock()
	defer p.Unlock()

	if _, ok := p.pools[poolID]; !ok {
		p.pools[poolID] = NewPool(p.timeout, p.idleTimeout, poolID)
	}

	// update pool size
	p.pools[poolID].SetSize(size)

	// Add the WebSocket connection to the pool
	p.pools[poolID].RegisterConnection(ws)

	return
}

func (p *Pools) Clean() {
	p.Lock()
	defer p.Unlock()

	if len(p.pools) == 0 {
		return
	}

	idle := 0
	busy := 0
	closed := 0

	for _, pool := range p.pools {
		if pool.IsEmpty() {
			log.Infof("Removing empty connection pool : %p", pool.id)
			pool.Shutdown()
			delete(p.pools, pool.id)
		}

		ps := pool.Size()
		idle += ps.Idle
		busy += ps.Busy
		closed = ps.Closed
	}

	log.Infof("%d pools, %d idle, %d busy, %d closed", len(p.pools), idle, busy, closed)
}

func (p *Pools) IsEmpty() bool {
	p.Lock()
	defer p.Unlock()
	return len(p.pools) == 0
}

func (p *Pools) MoreThenOne() bool {
	p.Lock()
	defer p.Unlock()
	return len(p.pools) > 1
}

func (p *Pools) GetPool(id PoolID) (pool *Pool, ok bool) {
	p.Lock()
	defer p.Unlock()
	pool, ok = p.pools[id]
	return
}
