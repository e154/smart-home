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

package modbus_rtu

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

type Actor struct {
	entity_manager.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      event_bus.EventBus
	actionPool    chan event_bus.EventCallAction
	stateMu       *sync.Mutex
}

func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus event_bus.EventBus) (actor *Actor) {

	actor = &Actor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService, adaptors),
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

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	actor.DeserializeAttr(entity.Attributes.Serialize())

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))
			a.ScriptEngine.PushFunction("ModbusRtu", NewModbusRtu(eventBus, actor))
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

func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *Actor) SetState(params entity_manager.EntityStateParams) error {

	oldState := e.GetEventState(e)

	e.Now(oldState)

	if params.NewState != nil {
		state := e.States[*params.NewState]
		e.State = &state
		e.State.ImageUrl = state.ImageUrl
	}

	e.AttrMu.Lock()
	e.Attrs.Deserialize(params.AttributeValues)
	e.AttrMu.Unlock()

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		Type:     e.Id.Type(),
		EntityId: e.Id,
		OldState: oldState,
		NewState: e.GetEventState(e),
	})

	return nil
}

func (e *Actor) addAction(event event_bus.EventCallAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg event_bus.EventCallAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId.Name(), action.Name); err != nil {
		log.Error(err.Error())
	}
}

func (e *Actor) localTopic(r string) string {
	var parent string
	if e.ParentId != nil {
		parent = e.ParentId.Name()
	}
	return fmt.Sprintf("%s/%s/%s", node.TopicPluginNode, parent, r)
}