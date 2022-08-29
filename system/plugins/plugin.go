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

package plugins

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/web"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

// Plugin ...
type Plugin struct {
	EntityManager entity_manager.EntityManager
	Adaptors      *adaptors.Adaptors
	ScriptService scripts.ScriptService
	PluginManager common.PluginManager
	EventBus      bus.Bus
	IsStarted     *atomic.Bool
	Scheduler     *scheduler.Scheduler
	Crawler       web.Crawler
}

// NewPlugin ...
func NewPlugin() *Plugin {
	return &Plugin{
		IsStarted: atomic.NewBool(false),
	}
}

// Load ...
func (p *Plugin) Load(service Service) error {
	p.Adaptors = service.Adaptors()
	p.EventBus = service.EventBus()
	p.EntityManager = service.EntityManager()
	p.ScriptService = service.ScriptService()
	p.PluginManager = service.PluginManager()
	p.Scheduler = service.Scheduler()
	p.Crawler = service.Crawler()

	if p.IsStarted.Load() {
		return ErrPluginIsLoaded
	}
	p.IsStarted.Store(true)

	return nil
}

// Unload ...
func (p *Plugin) Unload() error {

	if !p.IsStarted.Load() {
		return ErrPluginIsUnloaded
	}
	p.IsStarted.Store(false)

	return nil
}

// Name ...
func (p *Plugin) Name() string {
	panic("implement me")
}

// Type ...
func (p *Plugin) Type() PluginType {
	panic("implement me")
}

// Depends ...
func (p *Plugin) Depends() []string {
	panic("implement me")
}

// Version ...
func (p *Plugin) Version() string {
	panic("implement me")
}

// Options ...
func (p *Plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorCustomAttrs: false,
		ActorAttrs:       nil,
		ActorSetts:       nil,
	}
}

// LoadSettings ...
func (p *Plugin) LoadSettings(pl Plugable) (settings m.Attributes, err error) {
	var plugin m.Plugin
	if plugin, err = p.Adaptors.Plugin.GetByName(pl.Name()); err != nil {
		return
	}
	settings = plugin.Settings
	return
}

// GetPlugin ...
func (p *Plugin) GetPlugin(name string) (interface{}, error) {
	return p.PluginManager.GetPlugin(name)
}
