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
	"context"
	"fmt"

	"github.com/e154/smart-home/system/supervisor"

	"github.com/e154/smart-home/system/scheduler"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/logging"
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	pause   uint
	actor   *Actor
	entryId scheduler.EntryID
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin: supervisor.NewPlugin(),
		pause:  10,
	}
	return p
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}
	// every day at 00:00 am
	p.entryId, err = p.Scheduler.AddFunc("0 0 0 * * *", func() {
		p.actor.UpdateDay()
	})
	return p.load(service)
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}
	p.Scheduler.Remove(p.entryId)
	return p.unload()
}

// Load ...
func (p *plugin) load(service supervisor.Service) (err error) {

	var entity *m.Entity
	if entity, err = p.Adaptors.Entity.GetById(context.Background(), common.EntityId(fmt.Sprintf("%s.%s", EntityLogs, Name))); err == nil {

	}

	p.actor = NewActor(p.Supervisor, p.EventBus, entity)
	p.Supervisor.Spawn(p.actor.Spawn)

	logging.LogsHook = p.actor.LogsHook

	return
}

func (p *plugin) unload() (err error) {
	return
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	return p.load(nil)
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	return p.unload()
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
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
