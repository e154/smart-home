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

package logs

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/supervisor"

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
	actor   *Actor
	entryId scheduler.EntryID
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin: supervisor.NewPlugin(),
	}
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}
	// every day at 00:00 am
	p.entryId, err = p.Service.Scheduler().AddFunc("0 0 0 * * *", func() {
		p.actor.UpdateDay()
	})
	var entity *m.Entity
	if entity, err = p.Service.Adaptors().Entity.GetById(context.Background(), common.EntityId(fmt.Sprintf("%s.%s", EntityLogs, Name))); err != nil {
		entity = &m.Entity{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", EntityLogs, Name)),
			PluginName: "logs",
			Attributes: NewAttr(),
		}
		if err = p.Service.Adaptors().Entity.Add(context.Background(), entity); err != nil {
			return
		}
	}

	p.actor = NewActor(entity, service)
	p.AddActor(p.actor, entity)

	logging.LogsHook = p.actor.LogsHook

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	p.Service.Scheduler().Remove(p.entryId)
	err = p.Plugin.Unload(ctx)
	return
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
	return m.PluginOptions{}
}
