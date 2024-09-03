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

package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/e154/bus"
	"go.uber.org/fx"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/webpush"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = logger.MustGetLogger("handlers")
)

type EventHandler struct {
	stream   *stream.Stream
	eventBus bus.Bus
}

func NewEventHandler(lc fx.Lifecycle,
	stream *stream.Stream,
	eventBus bus.Bus) *EventHandler {
	handler := &EventHandler{
		stream:   stream,
		eventBus: eventBus,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			return handler.Start(ctx)
		},
		OnStop: func(ctx context.Context) (err error) {
			return handler.Shutdown(ctx)
		},
	})

	return handler
}

// Start ...
func (s *EventHandler) Start(_ context.Context) error {
	s.stream.Subscribe("event_get_last_state", s.EventGetLastState)
	s.stream.Subscribe("event_get_state_by_id", s.EventGetStateById)
	s.stream.Subscribe("event_add_webpush_subscription", s.EventAddWebPushSubscription)
	s.stream.Subscribe("event_get_webpush_public_key", s.EventGetWebPushPublicKey)
	s.stream.Subscribe("event_get_user_devices", s.EventGetUserDevices)
	s.stream.Subscribe("command_terminal", s.CommandTerminal)
	s.stream.Subscribe("event_get_server_version", s.EventGetServerVersion)
	s.stream.Subscribe("event_stt", s.EventSTT)
	return nil
}

// Shutdown ...
func (s *EventHandler) Shutdown(_ context.Context) error {
	s.stream.UnSubscribe("event_get_last_state")
	s.stream.UnSubscribe("event_get_state_by_id")
	s.stream.UnSubscribe("event_add_webpush_subscription")
	s.stream.UnSubscribe("event_get_webpush_public_key")
	s.stream.UnSubscribe("event_get_user_devices")
	s.stream.UnSubscribe("command_terminal")
	s.stream.UnSubscribe("event_get_server_version")
	s.stream.UnSubscribe("event_stt")
	return nil
}

func (s *EventHandler) EventGetWebPushPublicKey(client stream.IStreamClient, query string, body []byte) {
	var userID int64
	if user := client.GetUser(); user != nil {
		userID = user.Id
	}

	s.eventBus.Publish(webpush.TopicPluginWebpush, webpush.EventGetWebPushPublicKey{
		UserID:    userID,
		SessionID: client.SessionID(),
	})
}

func (s *EventHandler) EventGetUserDevices(client stream.IStreamClient, query string, body []byte) {
	var userID int64
	if user := client.GetUser(); user != nil {
		userID = user.Id
	}

	s.eventBus.Publish(webpush.TopicPluginWebpush, webpush.EventGetUserDevices{
		UserID:    userID,
		SessionID: client.SessionID(),
	})
}

func (s *EventHandler) EventAddWebPushSubscription(client stream.IStreamClient, query string, body []byte) {
	var userID int64
	if user := client.GetUser(); user != nil {
		userID = user.Id
	}

	subscription := &m.Subscription{}
	_ = json.Unmarshal(body, subscription)
	s.eventBus.Publish(webpush.TopicPluginWebpush, webpush.EventAddWebPushSubscription{
		UserID:       userID,
		SessionID:    client.SessionID(),
		Subscription: subscription,
	})
}

func (s *EventHandler) EventGetLastState(client stream.IStreamClient, query string, body []byte) {
	req := map[string]common.EntityId{}
	_ = json.Unmarshal(body, &req)
	id := req["entity_id"]
	s.eventBus.Publish("system/entities/"+id.String(), events.EventGetLastState{
		EntityId: id,
	})
}

func (s *EventHandler) EventGetStateById(client stream.IStreamClient, query string, body []byte) {
	req := events.EventGetStateById{}
	_ = json.Unmarshal(body, &req)
	s.eventBus.Publish("system/entities/"+req.EntityId.String(), events.EventGetStateById{
		EntityId:  req.EntityId,
		StorageId: req.StorageId,
		Common: events.Common{
			User:      client.GetUser(),
			SessionID: client.SessionID(),
		},
	})
}

func (s *EventHandler) CommandTerminal(client stream.IStreamClient, query string, body []byte) {
	s.eventBus.Publish("system/terminal", events.CommandTerminal{
		Common: events.Common{
			User:      client.GetUser(),
			SessionID: client.SessionID(),
		},
		Text: string(body),
	})
}

func (s *EventHandler) EventGetServerVersion(client stream.IStreamClient, query string, body []byte) {
	s.eventBus.Publish("system", events.EventGetServerVersion{
		Common: events.Common{
			User:      client.GetUser(),
			SessionID: client.SessionID(),
		},
	})
}

func (s *EventHandler) EventSTT(client stream.IStreamClient, query string, body []byte) {

	req := map[string]string{}
	_ = json.Unmarshal(body, &req)

	payload := req["payload"]
	payload = strings.Replace(payload, "data:audio/raw;base64,", "", 1)

	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		log.Error(err.Error())
		return
	}

	ioReaderData := bytes.NewReader(data)
	buf := bytes.NewBuffer(make([]byte, 0))
	if _, err := buf.ReadFrom(ioReaderData); err != nil {
		log.Error(err.Error())
		return
	}

	s.eventBus.Publish("system/stt", events.CommandSTT{
		Payload: buf,
	})
}
