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

package script

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

const (
	FuncEntityAction = "entityAction"
)

type EntityActor struct {
	entity_manager.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService *scripts.ScriptService
	system        entity_manager.IActorManager
	stateMu       *sync.Mutex
	actionPool    chan event_bus.EventCallAction
}

func NewEntityActor(entity *m.Entity,
	params map[string]interface{},
	adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) (actor *EntityActor, err error) {

	actor = &EntityActor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService),
		adaptors:      adaptors,
		scriptService: scriptService,
		actionPool:    make(chan event_bus.EventCallAction, 10),
		stateMu:       &sync.Mutex{},
	}

	actor.Attrs.Deserialize(params)

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))
			a.ScriptEngine.Do()
		}
	}

	// Script
	if actor.ScriptEngine != nil {
		actor.ScriptEngine.PushStruct("Actor", NewScriptBind(actor))
	}

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
	if changed, err := e.Attrs.Deserialize(params.AttributeValues); !changed {
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
