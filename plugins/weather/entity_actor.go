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

package weather

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	updateLock     *sync.Mutex
	positionLock   *sync.Mutex
	zoneAttributes m.EntityAttributes
	eventBus       event_bus.EventBus
}

func NewEntityActor(name string,
	eventBus event_bus.EventBus,
	entityManager entity_manager.EntityManager) *EntityActor {

	e := &EntityActor{
		BaseActor: entity_manager.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityWeather, name)),
			Name:              name,
			Description:       "weather plugin",
			EntityType:        EntityWeather,
			UnitOfMeasurement: "C°",
			AttrMu:            &sync.Mutex{},
			Attrs:             BaseForecast(),
			ParentId:          common.NewEntityId(fmt.Sprintf("%s.%s", zone.Name, name)),
			Manager:           entityManager,
		},
		eventBus:       eventBus,
		zoneAttributes: zone.NewAttr(),
		updateLock:     &sync.Mutex{},
		positionLock:   &sync.Mutex{},
	}

	return e
}

func (e *EntityActor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *EntityActor) setPosition(zoneAttr m.EntityAttributes) {
	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	changed, err := e.zoneAttributes.Deserialize(zoneAttr.Serialize())
	if err != nil {
		log.Error(err.Error())
		return
	}

	if !changed {
		return
	}

	e.eventBus.Publish(TopicPluginWeather, EventSubStateChanged{
		Type:       e.Id.Type(),
		EntityId:   e.Id,
		State:      SubStatePositionUpdate,
		Attributes: e.zoneAttributes.Copy(),
	})
}

func (e *EntityActor) SetState(params entity_manager.EntityStateParams) {

	log.Infof("update forecast for '%s'", e.Id)

	oldState := e.GetEventState(e)

	e.Now(oldState)

	e.AttrMu.Lock()
	e.Attrs.Deserialize(params.AttributeValues)
	e.AttrMu.Unlock()

	e.Send(entity_manager.MessageStateChanged{
		StorageSave: true,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})
}
