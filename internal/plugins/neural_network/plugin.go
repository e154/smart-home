// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package neural_network

import (
	"context"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

var (
	log = logger.MustGetLogger("plugins.neural_network")
)

var _ plugins.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
}

func New() plugins.Pluggable {
	return &plugin{
		Plugin: plugins.NewPlugin(),
	}
}

func (p *plugin) Load(ctx context.Context, service plugins.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)

	return nil
}

func (p *plugin) Unload(ctx context.Context) (err error) {
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)

	return nil
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor plugins.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	return
}

func (p *plugin) Name() string {
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
			actor.addAction(v)
		}
	}
}

func (p *plugin) Depends() []string {
	return nil
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:       true,
		ActorAttrs:   NewAttr(),
		ActorSetts:   NewSettings(),
		ActorActions: plugins.ToEntityActionShort(NewActions()),
		ActorStates:  plugins.ToEntityStateShort(NewStates()),
	}
}