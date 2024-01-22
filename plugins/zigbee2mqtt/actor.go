// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"context"
	"sync"

	"github.com/e154/smart-home/system/scripts"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	zigbee2mqttDevice *m.Zigbee2mqttDevice
	mqttMessageQueue  chan *Message
	actionPool        chan events.EventCallEntityAction
	newMsgMu          *sync.Mutex
	stateMu           *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor, err error) {

	var zigbee2mqttDevice *m.Zigbee2mqttDevice
	if zigbee2mqttDevice, err = service.Adaptors().Zigbee2mqttDevice.GetById(context.Background(), entity.Id.Name()); err != nil {
		return
	}

	actor = &Actor{
		BaseActor:         supervisor.NewBaseActor(entity, service),
		mqttMessageQueue:  make(chan *Message, 10),
		actionPool:        make(chan events.EventCallEntityAction, 1000),
		newMsgMu:          &sync.Mutex{},
		stateMu:           &sync.Mutex{},
		zigbee2mqttDevice: zigbee2mqttDevice,
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine.Engine() != nil {
			_, _ = a.ScriptEngine.Engine().Do()
		}
	}

	actor.ScriptsEngine.Spawn(func(engine *scripts.Engine) {
		engine.Do()
	})

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

}

func (e *Actor) Spawn() {

}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {

	e.SetActorState(params.NewState)
	e.DeserializeAttr(params.AttributeValues)
	e.SaveState(false, params.StorageSave)

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

	if _, err := e.ScriptsEngine.Engine().AssertFunction(FuncZigbee2mqttEvent, message); err != nil {
		log.Error(err.Error())
		return
	}
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	if action, ok := e.Actions[msg.ActionName]; ok {
		if action.ScriptEngine != nil && action.ScriptEngine.Engine() != nil {
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
				log.Error(err.Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.Engine().AssertFunction(FuncEntityAction, msg.EntityId, msg.ActionName, msg.Args); err != nil {
			log.Error(err.Error())
		}
	}
}
