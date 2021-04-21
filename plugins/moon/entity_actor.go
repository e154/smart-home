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

package moon

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/astronomics/moonphase"
	"github.com/e154/smart-home/common/astronomics/suncalc"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"math"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	positionLock        *sync.Mutex
	lat, lon, elevation float64
	moonAzimuth         float64
	moonElevation       float64
	phase               string
	horizonState        string
}

func NewEntityActor(name string) *EntityActor {

	entity := &EntityActor{
		BaseActor: entity_manager.BaseActor{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", EntityMoon, name)),
			Name:        name,
			Description: "moon plugin",
			EntityType:  EntityMoon,
			AttrMu:      &sync.Mutex{},
			Attrs:       NewAttr(),
		},
		positionLock: &sync.Mutex{},
	}

	return entity
}

func (e *EntityActor) Spawn(actorManager entity_manager.IActorManager) entity_manager.IActor {
	e.Manager = actorManager
	return e
}

func (e *EntityActor) setPosition(lat, lon, elevation float64, timezone int) {
	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	e.lat = lat
	e.lon = lon
	e.elevation = elevation
}

func (e *EntityActor) updateMoonPosition() {

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	oldState := e.GetEventState(e)

	now := e.Now(oldState)

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

	var attributeValues = make(m.EntityAttributeValue)
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

	log.Debugf("Moon horizonState %v", e.horizonState)

	e.DeserializeAttr(attributeValues)

	e.Send(entity_manager.MessageStateChanged{
		StorageSave: true,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})
}
