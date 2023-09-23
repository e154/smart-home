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
	"math"
	"sync"
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common/astronomics/moonphase"
	"github.com/e154/smart-home/common/astronomics/suncalc"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	positionLock  *sync.Mutex
	lat, lon      float64
	moonAzimuth   float64
	moonElevation float64
	phase         string
	horizonState  string
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor:    supervisor.NewBaseActor(entity, service),
		positionLock: &sync.Mutex{},
	}

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	actor.setPosition(entity.Settings)

	return actor
}

func (a *Actor) Destroy() {

}

func (e *Actor) Spawn() {
	e.UpdateMoonPosition(time.Now())
}

func (e *Actor) setPosition(settings m.Attributes) {
	if settings == nil {
		return
	}

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	if lat, ok := settings[AttrLat]; ok {
		e.lat = lat.Float64()
	}
	if lon, ok := settings[AttrLon]; ok {
		e.lon = lon.Float64()
	}
}

// UpdateMoonPosition ...
func (e *Actor) UpdateMoonPosition(now time.Time) {

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	oldState := e.GetEventState()

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

	//log.Debugf("Moon pos azimuth(%f), elevation(%f)", e.moonAzimuth, e.moonElevation)
	//log.Debugf("Moon phase %v", e.phase)

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

	//log.Debugf("Moon horizonState %v", e.horizonState)

	e.DeserializeAttr(attributeValues)

	e.Service.EventBus().Publish("system/entities/"+e.Id.String(), events.EventStateChanged{
		StorageSave: true,
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(),
	})
}
