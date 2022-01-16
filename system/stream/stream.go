// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package stream

import (
	"context"
	"go.uber.org/fx"
	"sync"

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
)

var (
	log = common.MustGetLogger("stream")
)

// Stream ...
type Stream struct {
	subMx       sync.Mutex
	subscribers map[string]func(client IStreamClient, id string, msg []byte)
	sesMx       sync.Mutex
	sessions    map[*Client]bool
}

// NewStreamService ...
func NewStreamService(lc fx.Lifecycle) *Stream {
	s := &Stream{
		subscribers: make(map[string]func(client IStreamClient, id string, msg []byte)),
		sessions:    make(map[*Client]bool),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return s.Start(ctx)
		},
		OnStop: func(ctx context.Context) (err error) {
			return s.Shutdown(ctx)
		},
	})

	return s
}

// Start ...
func (s *Stream) Start(cts context.Context) error {
	return nil
}

// Shutdown ...
func (s *Stream) Shutdown(cts context.Context) error {
	s.sesMx.Lock()
	defer s.sesMx.Unlock()

	for client, ok := range s.sessions {
		if !ok {
			continue
		}
		client.Close()
	}
	return nil
}

// Broadcast ...
func (s *Stream) Broadcast(query string, message []byte) {
	s.sesMx.Lock()
	defer s.sesMx.Unlock()

	for client, ok := range s.sessions {
		if !ok {
			continue
		}
		client.Broadcast(query, message)
	}
}

// Subscribe ...
func (s *Stream) Subscribe(command string, f func(IStreamClient, string, []byte)) {
	log.Infof("subscribe %s", command)
	s.subMx.Lock()
	defer s.subMx.Unlock()
	if s.subscribers[command] != nil {
		delete(s.subscribers, command)
	}
	s.subscribers[command] = f

}

// UnSubscribe ...
func (s *Stream) UnSubscribe(command string) {
	log.Infof("unsubscribe %s", command)
	s.subMx.Lock()
	defer s.subMx.Unlock()
	if s.subscribers[command] != nil {
		delete(s.subscribers, command)
	}
}

// NewConnection ...
func (s *Stream) NewConnection(server api.StreamService_SubscribeServer) error {

	client := NewClient(server)
	defer func() {
		log.Infof("websocket session closed")
		s.sesMx.Lock()
		delete(s.sessions, client)
		s.sesMx.Unlock()
	}()

	s.sesMx.Lock()
	s.sessions[client] = true
	s.sesMx.Unlock()

	log.Infof("new websocket session established")

	err := client.WritePump(s.Recv)
	return err
}

// Recv ...
func (s *Stream) Recv(client *Client, id, query string, b []byte) {
	s.subMx.Lock()
	f, ok := s.subscribers[query]
	s.subMx.Unlock()
	if ok {
		f(client, id, b)
	}
}
