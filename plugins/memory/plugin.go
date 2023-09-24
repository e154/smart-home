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

package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/e154/smart-home/system/supervisor"

	"github.com/e154/smart-home/common"

	m "github.com/e154/smart-home/models"
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	ticker *time.Ticker
	actor  *Actor
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

	if p.actor != nil {
		return
	}

	var entity *m.Entity
	if entity, err = p.Service.Adaptors().Entity.GetById(context.Background(), common.EntityId(fmt.Sprintf("%s.%s", EntityMemory, Name))); err != nil {
		entity = &m.Entity{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", EntityMemory, Name)),
			PluginName: Name,
			Metrics:    NewMetrics(),
			Attributes: NewAttr(),
		}
		err = p.Service.Adaptors().Entity.Add(context.Background(), entity)
	}

	p.actor = NewActor(entity, p.Service)
	p.AddActor(p.actor, entity)

	go func() {
		const pause = 10
		p.ticker = time.NewTicker(time.Second * time.Duration(pause))

		for range p.ticker.C {
			p.actor.selfUpdate()
		}
	}()

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}
	p.actor = nil
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
	return m.PluginOptions{
		ActorAttrs: NewAttr(),
	}
}
