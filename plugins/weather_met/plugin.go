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
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
)

const (
	// Name ...
	Name = "weather_met"
	// DefaultApiUrl ...
	DefaultApiUrl = "https://api.met.no/weatherapi/locationforecast/2.0/classic"
)

var (
	log = common.MustGetLogger("plugins.met")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	quit    chan struct{}
	pause   uint
	service common.PluginManager
	weather *WeatherMet
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin: plugins.NewPlugin(),
		quit:   make(chan struct{}),
		pause:  60,
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.weather = NewWeatherMet(p.EventBus, p.Adaptors)

	p.EventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)
	p.EventBus.Subscribe(weather.TopicPluginWeather, p.eventHandler)
	p.quit = make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Minute * time.Duration(p.pause))

		defer func() {
			ticker.Stop()
			close(p.quit)
		}()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				if err = p.weather.UpdateForecastForAll(); err != nil {
					log.Error(err.Error())
				}
			}
		}
	}()

	if err = p.weather.UpdateForecastForAll(); err != nil {
		log.Error(err.Error())
	}

	return
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.quit <- struct{}{}
	p.EventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)
	p.EventBus.Unsubscribe(weather.TopicPluginWeather, p.eventHandler)
	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(_ string, msg interface{}) {
	switch v := msg.(type) {
	case event_bus.EventAddedActor:
		if v.PluginName != "weather" {
			return
		}

		p.weather.AddWeather(v.EntityId, v.Settings)

	case weather.EventStateChanged:
		if v.Type != "weather" || v.State != weather.StatePositionUpdate {
			return
		}

		p.weather.UpdateWeatherList(v.EntityId, v.Settings)

	case event_bus.EventRemoveActor:
		if v.PluginName != "weather" {
			return
		}

		p.weather.RemoveWeather(v.EntityId)
	}
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return []string{weather.Name}
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}
