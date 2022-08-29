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

package weather_owm

import (
	"time"

	"github.com/e154/smart-home/common/events"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/weather"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/plugins"
)

const (
	// Name ...
	Name = "weather_owm"
)

var (
	log = logger.MustGetLogger("plugins.owm")
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
	weather *WeatherOwm
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

	// load settings
	var settings m.Attributes
	settings, err = p.LoadSettings(p)
	if err != nil {
		log.Warn(err.Error())
		settings = NewSettings()
	}

	if settings == nil {
		settings = NewSettings()
	}

	p.weather = NewWeatherOwm(p.EventBus, p.Adaptors, settings)

	_ = p.EventBus.Subscribe(bus.TopicEntities, p.eventHandler)
	_ = p.EventBus.Subscribe(weather.TopicPluginWeather, p.eventHandler)
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
	_ = p.EventBus.Unsubscribe(bus.TopicEntities, p.eventHandler)
	_ = p.EventBus.Unsubscribe(weather.TopicPluginWeather, p.eventHandler)
	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(_ string, msg interface{}) {
	switch v := msg.(type) {
	case events.EventAddedActor:
		if v.PluginName != "weather" {
			return
		}

		if name, ok := v.Settings[weather.AttrPlugin]; ok {
			if name.String() != Name {
				return
			}
			p.weather.AddWeather(v.EntityId, v.Settings)
		}

	case weather.EventStateChanged:
		if v.Type != "weather" || v.State != weather.StatePositionUpdate {
			return
		}

		p.weather.UpdateWeatherList(v.EntityId, v.Settings)

	case events.EventRemoveActor:
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
	return []string{}
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorSetts: NewSettings(),
	}
}
