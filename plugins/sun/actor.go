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

package sun

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/astronomics/suncalc"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"math"
	"sort"
	"sync"
)

type Actor struct {
	entity_manager.BaseActor
	positionLock        *sync.Mutex
	lat, lon, elevation float64
	solarAzimuth        float64
	solarElevation      float64
	phase               string
	horizonState        string
}

func NewActor(name string, entityManager entity_manager.EntityManager) *Actor {

	entity := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", EntitySun, name)),
			Name:        name,
			Description: "sun plugin",
			EntityType:  EntitySun,
			AttrMu:      &sync.Mutex{},
			Attrs:       NewAttr(),
			ParentId:    common.NewEntityId(fmt.Sprintf("%s.%s", zone.Name, name)),
			Manager:     entityManager,
			States:      States(),
		},
		positionLock: &sync.Mutex{},
	}

	return entity
}

func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (e *Actor) setPosition(lat, lon, elevation float64) {
	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	e.lat = lat
	e.lon = lon
	e.elevation = elevation
}

func (e *Actor) updateSunPosition() {

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	oldState := e.GetEventState(e)

	now := e.Now(oldState)

	sunrisePos := suncalc.GetSunPosition(now, e.lat, e.lon)
	e.solarAzimuth = sunrisePos.Azimuth*180/math.Pi + 180
	e.solarElevation = sunrisePos.Altitude * 180 / math.Pi

	times := suncalc.GetTimes(now, e.lat, e.lon)

	log.Debugf("Sun pos azimuth(%f), elevation(%f)", e.solarAzimuth, e.solarElevation)

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrAzimuth] = e.solarAzimuth
	attributeValues[AttrElevation] = e.solarElevation

	var dayTimes = make(DayTimes, 0)
	for _, t := range times {
		switch t.MorningName {
		case suncalc.Sunrise:
		case suncalc.Sunset:
		case suncalc.SunriseEnd:
		case suncalc.SunsetStart:
		case suncalc.Dawn:
		case suncalc.Dusk:
		case suncalc.NauticalDawn:
		case suncalc.NauticalDusk:
		case suncalc.NightEnd, suncalc.Night:
			continue
		case suncalc.GoldenHourEnd:
		case suncalc.GoldenHour:
		case suncalc.SolarNoon:
		case suncalc.Nadir:
		default:
			log.Warnf("unknown name %v", t.MorningName)
			continue
		}

		attributeValues[string(t.MorningName)] = t.Time

		dayTimes = append(dayTimes, DayTime{
			MorningName: string(t.MorningName),
			Time:        t.Time,
		})
	}
	sort.Sort(dayTimes)

	for _, t := range dayTimes {
		if now.Sub(t.Time).Minutes() > 0 {
			e.phase = t.MorningName
			if state, ok := e.States[t.MorningName]; ok {
				e.State = &state
			}
		}
	}
	log.Debugf("Sun phase %v", e.phase)

	attributeValues[AttrPhase] = e.phase

	if e.solarElevation > -0.833 {
		e.horizonState = StateAboveHorizon
	} else {
		e.horizonState = StateBelowHorizon
	}

	attributeValues[AttrHorizonState] = e.horizonState

	log.Debugf("Sun horizonState %v", e.horizonState)

	e.DeserializeAttr(attributeValues)

	e.Send(entity_manager.MessageStateChanged{
		StorageSave: true,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
	})
}
