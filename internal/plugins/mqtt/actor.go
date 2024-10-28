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
	"sync"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/events"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/mqtt"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/pkg/errors"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	message          *Message
	mqttMessageQueue chan *Message
	actionPool       chan events.EventCallEntityAction
	mqttClient       mqtt.MqttCli
	newMsgMu         *sync.Mutex
	stateMu          *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service,
	mqttClient mqtt.MqttCli) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:        supervisor.NewBaseActor(entity, service),
		message:          NewMessage(),
		mqttMessageQueue: make(chan *Message, 10),
		actionPool:       make(chan events.EventCallEntityAction, 1000),
		mqttClient:       mqttClient,
		newMsgMu:         &sync.Mutex{},
		stateMu:          &sync.Mutex{},
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushStruct("message", actor.message)
		}
	}

	if actor.ScriptsEngine != nil {
		actor.ScriptsEngine.PushStruct("message", actor.message)
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
	e.BaseActor.Spawn()
}

// SetState ...
func (e *Actor) SetState(params plugins.EntityStateParams) error {

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

	e.message.Update(message)

	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncMqttEvent, message); err != nil {
			log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
			return
		}
	}
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	if action, ok := e.Actions[msg.ActionName]; ok {
		if action.ScriptEngine != nil && action.ScriptEngine.Engine() != nil {
			if _, err := action.ScriptEngine.Engine().AssertFunction(FuncEntityAction, e.Id, action.Name, msg.Args); err != nil {
				log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
			}
			return
		}
	}
	if e.ScriptsEngine != nil && e.ScriptsEngine.Engine() != nil {
		if _, err := e.ScriptsEngine.AssertFunction(FuncEntityAction, e.Id, msg.ActionName, msg.Args); err != nil {
			log.Error(errors.Wrapf(err, "entity id: %s ", e.Id).Error())
		}
	}
}
