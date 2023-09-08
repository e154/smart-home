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
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/google/uuid"
	"go.uber.org/fx"

	"github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/webpush"
	"github.com/e154/smart-home/system/bus"
)

var (
	log = logger.MustGetLogger("stream")
)

// Stream ...
type Stream struct {
	*eventHandler
	subMx       sync.Mutex
	subscribers map[string]func(client IStreamClient, id string, msg []byte)
	sessions    sync.Map
	eventBus    bus.Bus
}

// NewStreamService ...
func NewStreamService(lc fx.Lifecycle,
	eventBus bus.Bus) (s *Stream) {
	s = &Stream{
		subscribers: make(map[string]func(client IStreamClient, id string, msg []byte)),
		sessions:    sync.Map{},
		eventBus:    eventBus,
	}

	s.eventHandler = NewEventHandler(s.Broadcast, s.DirectMessage)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return s.Start(ctx)
		},
		OnStop: func(ctx context.Context) (err error) {
			return s.Shutdown(ctx)
		},
	})

	return
}

// Start ...
func (s *Stream) Start(_ context.Context) error {
	_ = s.eventBus.Subscribe("system/entities/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Subscribe("system/plugins/html5_notify", s.eventHandler.eventHandler)
	_ = s.eventBus.Subscribe("system/automation/tasks/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Subscribe("system/automation/triggers/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Subscribe("system/automation/actions/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Subscribe(webpush.TopicPluginWebpush, s.eventHandler.eventHandler)
	_ = s.eventBus.Subscribe("system/services/mqtt", s.eventHandler.eventHandler)
	s.eventBus.Publish("system/services/stream", events.EventServiceStarted{Service: "Stream"})
	return nil
}

// Shutdown ...
func (s *Stream) Shutdown(_ context.Context) error {
	_ = s.eventBus.Unsubscribe("system/entities/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Unsubscribe("system/plugins/html5_notify", s.eventHandler.eventHandler)
	_ = s.eventBus.Unsubscribe("system/automation/tasks/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Unsubscribe("system/automation/triggers/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Unsubscribe("system/automation/actions/+", s.eventHandler.eventHandler)
	_ = s.eventBus.Unsubscribe("system/services/mqtt", s.eventHandler.eventHandler)
	_ = s.eventBus.Unsubscribe(webpush.TopicPluginWebpush, s.eventHandler.eventHandler)
	s.sessions.Range(func(key, value interface{}) bool {
		cli := value.(*Client)
		cli.Close()
		return true
	})
	s.eventBus.Publish("system/services/stream", events.EventServiceStopped{Service: "Stream"})
	return nil
}

// Broadcast ...
func (s *Stream) Broadcast(query string, message []byte) {
	s.sessions.Range(func(key, value interface{}) bool {
		cli := value.(*Client)
		_ = cli.Send(uuid.NewString(), query, message)
		return true
	})
}

// DirectMessage ...
func (s *Stream) DirectMessage(userID int64, query string, message []byte) {
	s.sessions.Range(func(key, value interface{}) bool {
		cli := value.(*Client)
		if userID != 0 && cli.user.Id != userID {
			return true
		}
		_ = cli.Send(uuid.NewString(), query, message)
		return true
	})
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
func (s *Stream) NewConnection(server api.StreamService_SubscribeServer, user *m.User) error {

	id := uuid.NewString()
	client := NewClient(server, user)
	defer func() {
		log.Infof("websocket session closed, email: '%s'", user.Email)
		s.sessions.Delete(id)
	}()

	s.sessions.Store(id, client)

	log.Infof("new websocket session established, email: '%s'", user.Email)

	err := client.WritePump(s.Recv)
	return err
}

// Recv ...
func (s *Stream) Recv(client *Client, id, query string, b []byte) {
	//log.Debugf("id: %s, query: %s, body: %s", id, query, string(b))
	s.subMx.Lock()
	f, ok := s.subscribers[query]
	s.subMx.Unlock()
	if ok {
		f(client, id, b)
	}
}
