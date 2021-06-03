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
	"errors"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
)

type Plugin struct {
	EntityManager entity_manager.EntityManager
	Adaptors      *adaptors.Adaptors
	ScriptService scripts.ScriptService
	EventBus      event_bus.EventBus
	IsStarted     *atomic.Bool
}

func NewPlugin() *Plugin {
	return &Plugin{
		IsStarted: atomic.NewBool(false),
	}
}

func (p Plugin) Load(service Service) error {
	p.Adaptors = service.Adaptors()
	p.EventBus = service.EventBus()
	p.EntityManager = service.EntityManager()
	p.ScriptService = service.ScriptService()

	if p.IsStarted.Load() {
		return errors.New("plugin is loaded")
	}
	p.IsStarted.Store(true)

	return nil
}

func (p Plugin) Unload() error {

	if !p.IsStarted.Load() {
		return errors.New("plugin is unloaded")
	}
	p.IsStarted.Store(false)

	return nil
}

func (p Plugin) Name() string {
	panic("implement me")
}

func (p Plugin) Type() PluginType {
	panic("implement me")
}

func (p Plugin) Depends() []string {
	panic("implement me")
}

func (p Plugin) Version() string {
	panic("implement me")
}

func (p Plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorCustomAttrs: false,
		ActorAttrs:       nil,
		ActorSetts:       nil,
	}
}
