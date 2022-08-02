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

package weather

import (
	"fmt"
	"sync"
	"time"

	"github.com/e154/smart-home/system/event_bus/events"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/astronomics/suncalc"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	updateLock   *sync.Mutex
	eventBus     event_bus.EventBus
	positionLock *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus) *Actor {

	name := entity.Id.Name()

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityWeather, name)),
			Name:              name,
			Description:       "weather plugin",
			EntityType:        EntityWeather,
			UnitOfMeasurement: "C°",
			AttrMu:            &sync.RWMutex{},
			ParentId:          common.NewEntityId(fmt.Sprintf("%s.%s", zone.Name, name)),
			Manager:           entityManager,
			States:            NewActorStates(false, false),
		},
		eventBus:     eventBus,
		updateLock:   &sync.Mutex{},
		positionLock: &sync.Mutex{},
	}

	if actor.Attrs == nil {
		actor.Attrs = BaseForecast()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	actor.UpdatePosition(entity.Settings)

	return actor
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {

	e.eventBus.Publish(TopicPluginWeather, EventStateChanged{
		Type:     e.Id.PluginName(),
		EntityId: e.Id,
		State:    StatePositionUpdate,
		Settings: e.Setts.Copy(),
	})
	return e
}

// UpdatePosition ...
func (e *Actor) UpdatePosition(settings m.Attributes) {
	if settings == nil {
		return
	}

	_, _ = e.Setts.Deserialize(settings.Serialize())

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	e.eventBus.Publish(TopicPluginWeather, EventStateChanged{
		Type:     e.Id.PluginName(),
		EntityId: e.Id,
		State:    StatePositionUpdate,
		Settings: e.Setts,
	})
}

// SetState ...
func (e *Actor) SetState(params entity_manager.EntityStateParams) error {

	log.Infof("update forecast for '%s'", e.Id)

	oldState := e.GetEventState(e)

	e.Now(oldState)

	// update sun position
	sunrisePos := suncalc.GetSunPosition(time.Now(), e.Setts[AttrLat].Float64(), e.Setts[AttrLon].Float64())
	var night = suncalc.IsNight(sunrisePos)
	var winter bool //todo fix

	if params.NewState != nil {
		state := e.States[*params.NewState]
		e.State = &state
		e.State.ImageUrl = GetImagePath(state.Name, night, winter)
	}

	e.AttrMu.Lock()
	_, _ = e.Attrs.Deserialize(params.AttributeValues)
	e.AttrMu.Unlock()

	e.eventBus.Publish(event_bus.TopicEntities, events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
		StorageSave: true,
	})

	return nil
}
