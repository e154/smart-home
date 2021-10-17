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

package scene

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
	FuncSceneEvent = "sceneEvent"
)

type Actor struct {
	entity_manager.BaseActor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	scriptEngine  *scripts.Engine
	eventPool     chan event_bus.EventCallScene
	stateMu       *sync.Mutex
}

func NewActor(entity *m.Entity,
	params map[string]interface{},
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	entityManager entity_manager.EntityManager) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventPool:     make(chan event_bus.EventCallScene, 10),
		stateMu:       &sync.Mutex{},
	}

	actor.Manager = entityManager
	actor.Attrs.Deserialize(params)

	// todo move to baseActor
	if len(entity.Scripts) != 0 {
		if actor.scriptEngine, err = scriptService.NewEngine(&entity.Scripts[0]); err != nil {
			return
		}
		actor.scriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.scriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
		actor.scriptEngine.Do()
	}

	// action worker
	go func() {
		for msg := range actor.eventPool {
			actor.runEvent(msg)
		}
	}()

	return
}

func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *Actor) addEvent(event event_bus.EventCallScene) {
	e.eventPool <- event
}

func (e *Actor) runEvent(msg event_bus.EventCallScene) {

	if _, err := e.scriptEngine.AssertFunction(FuncSceneEvent, msg.EntityId.Name()); err != nil {
		log.Error(err.Error())
	}
}
