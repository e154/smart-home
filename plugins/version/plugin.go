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

package version

import (
	"context"
	"fmt"
	"time"

	"github.com/e154/smart-home/common"

	"github.com/e154/smart-home/system/supervisor"

	m "github.com/e154/smart-home/models"
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	pause  uint
	ticker *time.Ticker
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
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	var entity *m.Entity
	if entity, err = p.Service.Adaptors().Entity.GetById(context.Background(), common.EntityId(fmt.Sprintf("%s.%s", EntityVersion, Name))); err != nil {
		entity = &m.Entity{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", EntityVersion, Name)),
			PluginName: Name,
			Attributes: NewAttr(),
		}
		err = p.Service.Adaptors().Entity.Add(context.Background(), entity)
	}

	go func() {
		p.ticker = time.NewTicker(time.Second * time.Duration(p.pause))

		for range p.ticker.C {
			p.Actors.Range(func(key, value any) bool {
				actor, _ := value.(*Actor)
				actor.selfUpdate()
				return true
			})
		}
	}()

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}
	p.Plugin.Unload(ctx)
	return nil
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
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
