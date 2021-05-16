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

package modbus_tcp

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      event_bus.EventBus
	actionPool    chan event_bus.EventCallAction
	stateMu       *sync.Mutex
}

func NewEntityActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus event_bus.EventBus) (actor *EntityActor) {

	actor = &EntityActor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		actionPool:    make(chan event_bus.EventCallAction, 10),
		stateMu:       &sync.Mutex{},
	}

	if actor.ParentId == nil {
		log.Warnf("entity %s, parent is nil", actor.Id)
	}

	actor.Manager = entityManager
	actor.Attrs = NewAttr()

	actor.DeserializeAttr(entity.Attributes.Serialize())

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))
			a.ScriptEngine.PushFunction("ModbusTcp",  NewModbusTcp(eventBus, actor))
			a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine == nil {
		return
	}

	// bind
	actor.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

func (e *EntityActor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *EntityActor) setState(params entity_manager.EntityStateParams) (changed bool) {
	return
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

func (e *EntityActor) localTopic(r string) string {
	var parent string
	if e.ParentId != nil {
		parent = e.ParentId.Name()
	}
	return fmt.Sprintf("%s/%s/%s", node.TopicPluginNode, parent, r)
}
