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

package stream

import (
	"encoding/json"

	"github.com/e154/smart-home/internal/plugins/webpush"
	"github.com/e154/smart-home/pkg/events"
)

type eventHandler struct {
	broadcast     func(query string, message []byte)
	directMessage func(userID int64, sessionID string, query string, message []byte)
}

func NewEventHandler(broadcast func(query string, message []byte),
	directMessage func(userID int64, sessionID string, query string, message []byte)) *eventHandler {
	return &eventHandler{
		broadcast:     broadcast,
		directMessage: directMessage,
	}
}

func (e *eventHandler) eventHandler(_ string, message interface{}) {

	switch v := message.(type) {

	// entities
	case events.EventStateChanged:
		go e.eventStateChangedHandler(message)
	case events.EventLastStateChanged:
		go e.eventStateChangedHandler(message)
	//case events.EventCreatedEntityModel:
	//case events.EventUpdatedEntityModel:
	case events.EventUpdatedMetric:
		go e.event(message)
	//case events.CommandUnloadEntity:
	case events.EventEntityLoaded:
		go e.event(message)
	case events.EventEntityUnloaded:
		go e.event(message)
	//case events.EventEntitySetState:
	case events.EventStateById:
		go e.eventStateById(v)

	// notifications
	case webpush.EventNewWebPushPublicKey:
		go e.eventNewWebPushPublicKey(v)
	case webpush.EventUserDevices:
		go e.eventUserDevices(v)
	case events.EventDirectMessage:
		go e.eventDirectMessage(v.UserID, v.SessionID, v.Query, v.Message)

	// plugins
	case events.EventPluginLoaded,
		events.EventPluginUnloaded:
		go e.event(message)

	// tasks
	case events.EventTaskLoaded,
		events.EventTaskUnloaded,
		events.EventTaskCompleted:
		go e.event(message)

	// triggers
	case events.EventTriggerLoaded,
		events.EventTriggerUnloaded,
		events.EventTriggerCompleted:
		go e.event(message)

	// actions
	case events.EventActionCompleted:
		go e.event(message)

	// mqtt
	case events.EventMqttNewClient:
		go e.event(message)

	// backup
	case events.EventCreatedBackup,
		events.EventRemovedBackup,
		events.EventUploadedBackup,
		events.EventStartedRestore:
		go e.event(message)

	// variables
	case events.EventRemovedVariableModel,
		events.EventUpdatedVariableModel:
		go e.event(message)

	// scripts
	case events.EventUpdatedScriptModel,
		events.EventRemovedScriptModel,
		events.EventCreatedScriptModel:
		go e.event(message)

	// version
	case events.EventServerVersion:
		go e.eventDirectMessage(v.UserId(), v.SessionID, events.EventName(v), v.Version)

	default:

	}
}

func (e *eventHandler) eventNewWebPushPublicKey(event webpush.EventNewWebPushPublicKey) {
	b, _ := json.Marshal(event)
	if event.UserID != 0 {
		e.directMessage(event.UserID, event.SessionID, events.EventName(event), b)
		return
	}
	e.broadcast(events.EventName(event), b)
}

func (e *eventHandler) eventStateById(event events.EventStateById) {
	b, _ := json.Marshal(event)
	if event.UserID != 0 {
		e.directMessage(event.UserID, event.SessionID, events.EventName(event), b)
		return
	}
	e.broadcast(events.EventName(event), b)
}

func (e *eventHandler) eventUserDevices(event webpush.EventUserDevices) {
	b, _ := json.Marshal(event)
	if event.UserID != 0 {
		e.directMessage(event.UserID, event.SessionID, events.EventName(event), b)
		return
	}
	e.broadcast(events.EventName(event), b)
}

func (e *eventHandler) eventStateChangedHandler(msg interface{}) {
	//todo optimize
	b, _ := json.Marshal(msg)
	e.broadcast("state_changed", b)
}

func (e *eventHandler) eventDirectMessage(userID int64, sessionID string, query string, msg interface{}) {
	b, _ := json.Marshal(msg)
	e.directMessage(userID, sessionID, query, b)
}

func (e *eventHandler) event(msg interface{}) {
	b, _ := json.Marshal(msg)
	e.broadcast(events.EventName(msg), b)
}
