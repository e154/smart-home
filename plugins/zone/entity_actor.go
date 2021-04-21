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

package zone

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	entities []entity_manager.IActor
	stateMu  *sync.Mutex
}

func NewEntityActor(name string, params m.EntityAttributeValue) *EntityActor {

	attributes := NewAttr()
	attributes.Deserialize(params)

	e := &EntityActor{
		BaseActor: entity_manager.BaseActor{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", EntityZone, name)),
			Name:       name,
			EntityType: EntityZone,
			AttrMu:     &sync.Mutex{},
			Attrs:      attributes,
		},
		stateMu: &sync.Mutex{},
	}

	return e
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

	var changed bool
	var err error
	e.AttrMu.Lock()
	if changed, err = e.Attrs.Deserialize(params.AttributeValues); !changed {
		if err != nil {
			log.Warn(err.Error())
		}

		if oldState.LastUpdated != nil {
			delta := now.Sub(*oldState.LastUpdated).Milliseconds()
			if delta < 200 {
				e.AttrMu.Unlock()
				return
			}
		}
	}
	e.AttrMu.Unlock()

	go e.Send(entity_manager.MessageStateChanged{
		StorageSave: true,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})
}
