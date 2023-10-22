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

package mqtt

import (
	"fmt"
	"sync"

	"github.com/e154/smart-home/system/scripts"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	message          *Message
	mqttMessageQueue chan *Message
	actionPool       chan events.EventCallEntityAction
	mqttClient       mqtt.MqttCli
	newMsgMu         *sync.Mutex
	stateMu          *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service,
	mqttClient mqtt.MqttCli) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:        supervisor.NewBaseActor(entity, service),
		message:          NewMessage(),
		mqttMessageQueue: make(chan *Message, 10),
		actionPool:       make(chan events.EventCallEntityAction, 10),
		mqttClient:       mqttClient,
		newMsgMu:         &sync.Mutex{},
		stateMu:          &sync.Mutex{},
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine.Engine() != nil {
			_, _ = a.ScriptEngine.Engine().Do()
		}
	}

	for _, engine := range actor.ScriptEngines {
		engine.Spawn(func(engine *scripts.Engine) {
			engine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			engine.PushStruct("message", actor.message)
			engine.Do()
		})
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	// mqtt worker
	go func() {
		for message := range actor.mqttMessageQueue {
			actor.mqttNewMessage(message)
		}
	}()

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return
}

func (e *Actor) Destroy() {
	if e.Setts != nil && e.Setts[AttrSubscribeTopic] != nil {
		e.mqttClient.Unsubscribe(e.Setts[AttrSubscribeTopic].String())
	}
}

// Spawn ...
func (e *Actor) Spawn() {

	if e.Setts != nil && e.Setts[AttrSubscribeTopic] != nil {
		_ = e.mqttClient.Subscribe(e.Setts[AttrSubscribeTopic].String(), e.mqttOnPublish)
	}

	return
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	oldState := e.GetEventState()
	now := e.Now(oldState)

	if params.NewState != nil {
		if state, ok := e.States[*params.NewState]; ok {
			e.State = &state
		}
	}

	e.AttrMu.Lock()
	changed, err := e.Attrs.Deserialize(params.AttributeValues)
	if !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				e.AttrMu.Unlock()
				return nil
			}
		}
	}
	e.AttrMu.Unlock()

	go e.SaveState(events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
		StorageSave: params.StorageSave,
	})

	return nil
}

func (e *Actor) mqttOnPublish(client mqtt.MqttCli, msg mqtt.Message) {
	message := NewMessage()
	message.Payload = string(msg.Payload)
	message.Topic = msg.Topic
	message.Qos = msg.Qos
	message.Duplicate = msg.Dup

	e.mqttMessageQueue <- message
}

func (e *Actor) mqttNewMessage(message *Message) {

	e.newMsgMu.Lock()
	defer e.newMsgMu.Unlock()

	e.message.Update(message)
	for _, engine := range e.ScriptEngines {
		if _, err := engine.Engine().AssertFunction(FuncMqttEvent, message); err != nil {
			log.Error(err.Error())
			return
		}
	}
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine.Engine() == nil {
		return
	}
	if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
		log.Error(err.Error())
	}
}
