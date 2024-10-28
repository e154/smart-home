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

package cpuspeed

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

var _ plugins.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	ticker *time.Ticker
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
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	if _, err = p.Service.Adaptors().Entity.GetById(context.Background(), common.EntityId(fmt.Sprintf("%s.%s", EntityCpuspeed, Name))); err != nil {
		entity := &m.Entity{
			Id:          common.EntityId(fmt.Sprintf("%s.%s", EntityCpuspeed, Name)),
			Description: "cpu usage",
			PluginName:  Name,
			Metrics:     NewMetrics(),
			Attributes:  NewAttr(),
		}
		err = p.Service.Adaptors().Entity.Add(context.Background(), entity)
	}

	go func() {
		const pause = 10
		p.ticker = time.NewTicker(time.Second * time.Duration(pause))

		for range p.ticker.C {
			p.Actors.Range(func(key, value any) bool {
				actor, _ := value.(*Actor)
				actor.selfUpdate()
				return true
			})
		}
	}()

	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}
	err = p.Plugin.Unload(ctx)
	return
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor plugins.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	if entity.Metrics == nil {
		entity.Metrics = NewMetrics()
	}
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
	return m.PluginOptions{
		ActorAttrs: NewAttr(),
	}
}
