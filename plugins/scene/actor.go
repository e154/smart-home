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
	"sync"

	"github.com/e154/smart-home/common/events"

	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	eventPool chan events.EventCallScene
	stateMu   *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) (actor *Actor, err error) {

	actor = &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		eventPool: make(chan events.EventCallScene, 99),
		stateMu:   &sync.Mutex{},
	}

	// action worker
	go func() {
		for msg := range actor.eventPool {
			actor.runEvent(msg)
		}
	}()

	return
}

func (e *Actor) Destroy() {
	close(e.eventPool)
}

func (e *Actor) Spawn() {

}

func (e *Actor) addEvent(event events.EventCallScene) {
	e.eventPool <- event
}

func (e *Actor) runEvent(msg events.EventCallScene) {
	if _, err := e.ScriptEngine.AssertFunction(FuncSceneEvent, msg.EntityId); err != nil {
		log.Error(err.Error())
	}
}
