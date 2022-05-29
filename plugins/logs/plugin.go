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

package logs

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/plugins"
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	quit  chan struct{}
	pause uint
	actor *Actor
	cron  *cron.Cron
	task  *cron.Task
}

// New ...
func New() plugins.Plugable {
	p := &plugin{
		Plugin: plugins.NewPlugin(),
		pause:  10,
		cron: cron.NewCron(),
	}
	return p
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.actor = NewActor(p.EntityManager, p.EventBus)
	p.EntityManager.Spawn(p.actor.Spawn)

	logging.LogsHook = p.actor.LogsHook

	// Spawn ...
	p.task, _ = p.cron.NewTask("0 0 0 * * *", func() {
		p.actor.UpdateDay()
	})
	p.cron.Run()

	return nil
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	if p.task != nil {
		p.cron.RemoveTask(p.task)
	}
	p.cron.Stop()
	p.task = nil
	p.quit <- struct{}{}
	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorAttrs: NewAttr(),
	}
}
