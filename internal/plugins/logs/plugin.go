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
	"embed"
	"fmt"

	"github.com/e154/smart-home/internal/system/logging"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scheduler"
)

var _ plugins.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

type plugin struct {
	*plugins.Plugin
	actor   *Actor
	entryId scheduler.EntryID
}

// New ...
func New() plugins.Pluggable {
	p := &plugin{
		Plugin: plugins.NewPlugin(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service plugins.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}
	// every day at 00:00 am
	p.entryId, _ = p.Service.Scheduler().AddFunc("0 0 0 * * *", func() {
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
	_ = p.AddActor(p.actor, entity)

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

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{}
}
