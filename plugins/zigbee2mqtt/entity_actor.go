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

package zigbee2mqtt

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	adaptors          *adaptors.Adaptors
	scriptService     *scripts.ScriptService
	zigbee2mqttDevice *m.Zigbee2mqttDevice
	message           *Message
	mqttMessageQueue  chan *Message
	actionPool        chan event_bus.EventCallAction
	newMsgMu          *sync.Mutex
	stateMu           *sync.Mutex
}

func NewEntityActor(entity *m.Entity,
	params map[string]interface{},
	adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) (actor *EntityActor, err error) {

	var zigbee2mqttDevice *m.Zigbee2mqttDevice
	if zigbee2mqttDevice, err = adaptors.Zigbee2mqttDevice.GetById(entity.Id.Name()); err != nil {
		return nil, err
	}

	actor = &EntityActor{
		BaseActor:         entity_manager.NewBaseActor(entity, scriptService),
		adaptors:          adaptors,
		scriptService:     scriptService,
		zigbee2mqttDevice: zigbee2mqttDevice,
		message:           NewMessage(),
		mqttMessageQueue:  make(chan *Message, 10),
		actionPool:        make(chan event_bus.EventCallAction, 10),
		newMsgMu:          &sync.Mutex{},
		stateMu:           &sync.Mutex{},
	}

	actor.Attrs.Deserialize(params)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))
			a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine == nil {
		return
	}

	// message
	actor.ScriptEngine.PushStruct("message", actor.message)

	// bind
	actor.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))

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

func (e *EntityActor) Spawn(actorManager entity_manager.IActorManager) entity_manager.IActor {
	e.Manager = actorManager
	return e
}

func (e *EntityActor) SetState(params entity_manager.EntityStateParams) {
	if !e.setState(params) {
		return
	}

	message := NewMessage()
	message.NewState = params

	e.mqttMessageQueue <- message
}

func (e *EntityActor) setState(params entity_manager.EntityStateParams) (changed bool) {
	e.stateMu.Lock()
	defer e.stateMu.Unlock()

	oldState := e.GetEventState(e)
	now := e.Now(oldState)

	if params.NewState != nil {
		if state, ok := e.States[*params.NewState]; ok {
			e.State = &state
		}
	}

	e.AttrMu.Lock()
	var err error
	if changed, err = e.Attrs.Deserialize(params.AttributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			//fmt.Println("delta", delta)
			if delta < 200 {
				e.AttrMu.Unlock()
				return
			}
		}
	}
	e.AttrMu.Unlock()

	e.Send(entity_manager.MessageStateChanged{
		StorageSave: true,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})

	return
}

func (e *EntityActor) mqttOnPublish(client *mqtt.Client, msg mqtt.Message) {
	message := NewMessage()
	message.Payload = string(msg.Payload)
	message.Topic = msg.Topic
	message.Qos = msg.Qos
	message.Duplicate = msg.Dup

	e.mqttMessageQueue <- message
}

func (e *EntityActor) mqttNewMessage(message *Message) {
	e.newMsgMu.Lock()
	defer e.newMsgMu.Unlock()

	e.message.Update(message)
	if e.ScriptEngine == nil {
		return
	}
	if _, err := e.ScriptEngine.AssertFunction(FuncZigbee2mqttEvent); err != nil {
		log.Error(err.Error())
		return
	}
}

func (e *EntityActor) addAction(event event_bus.EventCallAction) {
	e.actionPool <- event
}

func (e *EntityActor) runAction(msg event_bus.EventCallAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId.Name(), action.Name); err != nil {
		log.Error(err.Error())
	}
}
