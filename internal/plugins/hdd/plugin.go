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

package hdd

import (
	"context"
	"embed"
	"fmt"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/events"
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

	if list, _ := p.Service.Adaptors().Entity.GetByType(ctx, Name, 5, 0); len(list) == 0 {
		entity := &m.Entity{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", EntityHDD, Name)),
			PluginName: Name,
			Attributes: NewAttr(),
			Metrics:    NewMetrics(),
		}
		_ = p.Service.Adaptors().Entity.Add(context.Background(), entity)
	}

	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)
	return
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)
	err = p.Plugin.Unload(ctx)
	return
}

// Name ...
func (p plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(topic string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventStateChanged:
	case events.EventCallEntityAction:
		values, ok := p.Check(v)
		if !ok {
			return
		}
		for _, value := range values {
			actor := value.(*Actor)
			actor.runAction(v)
		}
	}
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor plugins.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	if entity.Metrics == nil {
		entity.Metrics = NewMetrics()
	}
	return
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorAttrs:   NewAttr(),
		ActorSetts:   NewSettings(),
		ActorActions: plugins.ToEntityActionShort(NewActions()),
	}
}
