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

package script

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

const (
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus      event_bus.EventBus
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	system        entity_manager.EntityManager
	stateMu       *sync.Mutex
	actionPool    chan event_bus.EventCallAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		actionPool:    make(chan event_bus.EventCallAction, 10),
		stateMu:       &sync.Mutex{},
		eventBus:      eventBus,
	}

	actor.Manager = entityManager

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
			a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			a.ScriptEngine.Do()
		}
	}

	// Script
	if actor.ScriptEngine != nil {
		actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

// SetState ...
func (e *Actor) SetState(params entity_manager.EntityStateParams) (err error) {

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
	var changed bool
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

	e.eventBus.Publish(event_bus.TopicEntities, event_bus.EventStateChanged{
		Type:     e.Id.Type(),
		EntityId: e.Id,
		OldState: oldState,
		NewState: e.GetEventState(e),
	})

	return
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
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name); err != nil {
		log.Error(err.Error())
	}
}
