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

package weather_met

import (
	"fmt"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/astronomics/suncalc"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/scripts"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	Zone
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	eventBus      bus.Bus
	actionPool    chan events.EventCallAction
	weather       *WeatherMet
	winter        bool
	theme         string
}

// NewActor ...
func NewActor(entity *m.Entity,
	entityManager entity_manager.EntityManager,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus bus.Bus,
	crawler web.Crawler) *Actor {

	actor := &Actor{
		BaseActor:     entity_manager.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		actionPool:    make(chan events.EventCallAction, 10),
		weather:       NewWeatherMet(adaptors, crawler),
	}

	actor.Manager = entityManager

	if actor.Attrs == nil {
		actor.Attrs = weather.BaseForecast()
	}

	if actor.Setts == nil {
		actor.Setts = weather.NewSettings()
	}

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", entity_manager.NewScriptBind(actor))
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

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

func (e *Actor) destroy() {

}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	go e.updateForecast()
	return e
}

// SetState ...
func (e *Actor) SetState(params entity_manager.EntityStateParams) error {
	return nil
}

func (e *Actor) addAction(event events.EventCallAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	if action.ScriptEngine == nil {
		return
	}
	if _, err := action.ScriptEngine.AssertFunction(FuncEntityAction, msg.EntityId, action.Name); err != nil {
		log.Error(err.Error())
	}
}

func (e *Actor) update() {

	oldState := e.GetEventState(e)

	e.Now(oldState)

	if !e.updateForecast() {
		return
	}

	e.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		PluginName:  e.Id.PluginName(),
		EntityId:    e.Id,
		OldState:    oldState,
		NewState:    e.GetEventState(e),
		StorageSave: true,
	})

}

func (e *Actor) updateForecast() (changed bool) {
	forecast, err := e.weather.UpdateForecast(e.Zone)
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
		state := e.States[val.String()]
		e.State = &state
		e.State.ImageUrl = weather.GetImagePath(state.Name, e.theme, night, e.winter)
	}

	for i := 1; i < 6; i++ {
		if main, ok := e.Attrs[fmt.Sprintf("day%d_main", i)]; ok {
			e.Attrs[fmt.Sprintf("day%d_icon", i)].Value = weather.GetImagePath(main.String(), e.theme, night, e.winter)
		}
	}

	return
}
