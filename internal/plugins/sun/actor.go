// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	"math"
	"sort"
	"sync"
	"time"

	"github.com/e154/smart-home/internal/common/astronomics/suncalc"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	positionLock        *sync.Mutex
	lat, lon, elevation float64
	solarAzimuth        float64
	solarElevation      float64
	phase               string
	horizonState        string
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) *Actor {

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

func (e *Actor) Destroy() {

}

// Spawn ...
func (e *Actor) Spawn() {
	e.UpdateSunPosition(time.Now())
}

func (e *Actor) setPosition(settings m.Attributes) {
	if settings == nil {
		return
	}

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	if settings[AttrLat] != nil {
		e.lat = settings[AttrLat].Float64()
	}
	if settings[AttrLon] != nil {
		e.lon = settings[AttrLon].Float64()
	}
}

// UpdateSunPosition ...
func (e *Actor) UpdateSunPosition(now time.Time) {

	e.positionLock.Lock()
	defer e.positionLock.Unlock()

	sunrisePos := suncalc.GetSunPosition(now, e.lat, e.lon)
	e.solarAzimuth = sunrisePos.Azimuth*180/math.Pi + 180
	e.solarElevation = sunrisePos.Altitude * 180 / math.Pi

	times := suncalc.GetTimes(now, e.lat, e.lon)

	//log.Debugf("Sun pos azimuth(%f), elevation(%f)", e.solarAzimuth, e.solarElevation)

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
			e.SetActorState(common.String(t.MorningName))
		}
	}
	//log.Debugf("Sun phase %v", e.phase)

	attributeValues[AttrPhase] = e.phase

	if e.solarElevation > -0.833 {
		e.horizonState = StateAboveHorizon
	} else {
		e.horizonState = StateBelowHorizon
	}

	attributeValues[AttrHorizonState] = e.horizonState

	//log.Debugf("Sun horizonState %v", e.horizonState)

	e.DeserializeAttr(attributeValues)

	e.SaveState(false, true)
}
