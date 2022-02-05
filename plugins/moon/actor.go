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

package moon

import (
	"fmt"
	"github.com/e154/smart-home/system/event_bus/events"
	"math"
	"sync"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/astronomics/moonphase"
	"github.com/e154/smart-home/common/astronomics/suncalc"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	eventBus      event_bus.EventBus
	positionLock  *sync.Mutex
	lat, lon      float64
	moonAzimuth   float64
	moonElevation float64
	phase         string
	horizonState  string
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus) *Actor {

	name := entity.Id.Name()

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", EntityMoon, name)),
			Name:        name,
			Description: "moon plugin",
			EntityType:  EntityMoon,
			AttrMu:      &sync.RWMutex{},
			Attrs:       entity.Attributes,
			Setts:       entity.Settings,
			Manager:     entityManager,
			States:      NewStates(),
		},
		eventBus:     eventBus,
		positionLock: &sync.Mutex{},
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	actor.setPosition(entity.Settings)

	return actor
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	e.UpdateMoonPosition(time.Now())
	return e
}

func (e *Actor) setPosition(settings m.Attributes) {
	if settings == nil {
		return
	}

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	e.lat = settings[AttrLat].Float64()
	e.lon = settings[AttrLon].Float64()
}

// UpdateMoonPosition ...
func (e *Actor) UpdateMoonPosition(now time.Time) {

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	oldState := e.GetEventState(e)

	e.Now(oldState)

	moon := moonphase.New(now)
	//fmt.Println(moon.PhaseName())

	phase := int(math.Floor((moon.Phase() + 0.0625) * 8))
	switch phase {
	case 0:
		e.phase = StateNewMoon
	case 1:
		e.phase = StateWaxingCrescent
	case 2:
		e.phase = StateFirstQuarter
	case 3:
		e.phase = StateWaxingGibbous
	case 4:
		e.phase = StateFullMoon
	case 5:
		e.phase = StateWaningGibbous
	case 6:
		e.phase = StateThirdQuarter
	case 7:
		e.phase = StateWaningCrescent
	case 8:
		e.phase = StateNewMoon
	default:
		log.Warnf("unknown status %s", moon.PhaseName())
	}

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrPhase] = e.phase

	moonPosition := suncalc.GetMoonPosition(now, e.lat, e.lon)
	e.moonAzimuth = moonPosition.Azimuth*180/math.Pi + 180
	e.moonElevation = moonPosition.Altitude * 180 / math.Pi

	log.Debugf("Moon pos azimuth(%f), elevation(%f)", e.moonAzimuth, e.moonElevation)
	log.Debugf("Moon phase %v", e.phase)

	attributeValues[AttrAzimuth] = e.moonAzimuth
	attributeValues[AttrElevation] = e.moonElevation

	if e.moonElevation > -0.833 {
		e.horizonState = StateAboveHorizon
	} else {
		e.horizonState = StateBelowHorizon
	}

	attributeValues[AttrHorizonState] = e.horizonState

	if state, ok := e.States[e.horizonState]; ok {
		e.State = &state
	}

	log.Debugf("Moon horizonState %v", e.horizonState)

	e.DeserializeAttr(attributeValues)

	e.eventBus.Publish(event_bus.TopicEntities, events.EventStateChanged{
		StorageSave: true,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})
}
