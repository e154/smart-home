// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package graceful_service

import (
	"sync"
)

type IGracefulClient interface {
	Shutdown()
}

type GracefulServicePool struct {
	cfg     *GracefulServiceConfig
	m       sync.Mutex
	clients map[int]IGracefulClient
}

func NewGracefulServicePool(cfg *GracefulServiceConfig) *GracefulServicePool {
	return &GracefulServicePool{
		cfg:     cfg,
		clients: make(map[int]IGracefulClient),
	}
}

func (h *GracefulServicePool) subscribe(client IGracefulClient) (id int) {
	h.m.Lock()
	id = len(h.clients)
	h.clients[id] = client
	h.m.Unlock()
	return
}

func (h *GracefulServicePool) unsubscribe(id int) {
	h.m.Lock()
	if _, ok := h.clients[id]; ok {
		delete(h.clients, id)
	}
	h.m.Unlock()
}

func (h *GracefulServicePool) shutdown() {
	h.m.Lock()
	i := len(h.clients)
	for ;i>=0;i-- {
		client := h.clients[i]
		if client != nil {
			client.Shutdown()
		}
	}
	h.m.Unlock()
}
