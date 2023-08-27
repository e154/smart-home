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
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/node"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	actionPool    chan events.EventCallEntityAction
	stateMu       *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	visor supervisor.Supervisor,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus bus.Bus) (actor *Actor) {

	actor = &Actor{
		BaseActor:     supervisor.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		actionPool:    make(chan events.EventCallEntityAction, 10),
		stateMu:       &sync.Mutex{},
	}

	//if actor.ParentId == nil {
	//	log.Warnf("entity %s, parent is nil", actor.Id)
	//}

	actor.Supervisor = visor

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
			a.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			a.ScriptEngine.PushFunction("ModbusTcp", NewModbusTcp(eventBus, actor))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		actor.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

// Spawn ...
func (e *Actor) Spawn() supervisor.PluginActor {
	return e
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {

	oldState := e.GetEventState(e)

	e.Now(oldState)

	if params.NewState != nil {
		state := e.States[*params.NewState]
		e.State = &state
		e.State.ImageUrl = state.ImageUrl
	}

	e.AttrMu.Lock()
	_, _ = e.Attrs.Deserialize(params.AttributeValues)
	e.AttrMu.Unlock()

	e.eventBus.Publish("system/entities/"+e.Id.String(), events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
		StorageSave: params.StorageSave,
	})

	return nil
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
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name); err != nil {
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
