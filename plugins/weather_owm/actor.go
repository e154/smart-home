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

package weather_owm

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"time"

	"github.com/e154/smart-home/common/astronomics/suncalc"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	Zone
	weather *WeatherOwm
	winter  bool
	theme   string
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
		weather:   NewWeatherOwm(service.Adaptors(), service.Crawler()),
	}

	if actor.Attrs == nil {
		actor.Attrs = weather.BaseForecast()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	// zone
	if lon, ok := actor.Setts[weather.AttrLon]; ok {
		actor.Zone.Lon = lon.Float64()
	}
	actor.Zone.Name = actor.Id.Name()
	if lat, ok := actor.Setts[weather.AttrLat]; ok {
		actor.Zone.Lat = lat.Float64()
	}
	actor.Zone.Name = entity.Id.Name()

	// etc
	if theme, ok := actor.Setts[weather.AttrTheme]; ok {
		actor.theme = theme.String()
	}
	if winter, ok := actor.Setts[weather.AttrWinter]; ok {
		actor.winter = winter.Bool()
	}

	return actor
}

func (e *Actor) Destroy() {

}

// Spawn ...
func (e *Actor) Spawn() {
	go e.updateForecast()
}

// SetState ...
func (e *Actor) SetState(params supervisor.EntityStateParams) error {
	return nil
}

func (e *Actor) update() {

	if !e.updateForecast() {
		return
	}

	e.SaveState(false, true)

}

func (e *Actor) updateForecast() (changed bool) {
	forecast, err := e.weather.UpdateForecast(e.Zone, e.Setts)
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Infof("update forecast for '%s'", e.Id)

	e.AttrMu.Lock()
	defer e.AttrMu.Unlock()
	if changed, err = e.Attrs.Deserialize(forecast); err != nil {
		return
	}

	// update sun position
	sunrisePos := suncalc.GetSunPosition(time.Now(), e.Setts[weather.AttrLat].Float64(), e.Setts[weather.AttrLon].Float64())
	var night = suncalc.IsNight(sunrisePos)

	if val, ok := e.Attrs[weather.AttrWeatherMain]; ok {
		e.SetActorState(common.String(val.String()))
		e.SetActorStateImage(weather.GetImagePath(val.String(), e.theme, night, e.winter), nil)
	}

	for i := 1; i < 6; i++ {
		if main, ok := e.Attrs[fmt.Sprintf("day%d_main", i)]; ok {
			e.Attrs[fmt.Sprintf("day%d_icon", i)].Value = weather.GetImagePath(main.String(), e.theme, night, e.winter)
		}
	}

	return
}
