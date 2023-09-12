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
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	eventBus      bus.Bus
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	system        supervisor.Supervisor
	stateMu       *sync.Mutex
	actionPool    chan events.EventCallEntityAction
}

// NewActor ...
func NewActor(entity *m.Entity,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	visor supervisor.Supervisor,
	eventBus bus.Bus) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:     supervisor.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		actionPool:    make(chan events.EventCallEntityAction, 10),
		stateMu:       &sync.Mutex{},
		eventBus:      eventBus,
	}

	actor.Supervisor = visor

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			a.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			_, _ = a.ScriptEngine.Do()
		}
	}

	// Script
	if actor.ScriptEngine != nil {
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
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
func (e *Actor) Spawn() supervisor.PluginActor {
	return e
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) (err error) {

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

	e.eventBus.Publish("system/entities/"+e.Id.String(), events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
		StorageSave: params.StorageSave,
	})

	return
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
	if action.ScriptEngine == nil {
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name, msg.Args); err != nil {
		log.Error(err.Error())
	}
}
