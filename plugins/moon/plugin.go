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

package moon

import (
	"context"
	"time"

	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.moon")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	ticker *time.Ticker
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin: supervisor.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	go func() {
		const pause = 240
		p.ticker = time.NewTicker(time.Second * pause)

		for range p.ticker.C {
			p.updatePositionForAll()
		}
	}()

	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	p.ticker.Stop()
	err = p.Plugin.Unload(ctx)
	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	return
}

func (p *plugin) updatePositionForAll() {
	now := time.Now()
	p.Actors.Range(func(key, value any) bool {
		actor := value.(*Actor)
		actor.UpdateMoonPosition(now)
		return true
	})
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
		Actors:      true,
		ActorAttrs:  NewAttr(),
		ActorSetts:  NewSettings(),
		ActorStates: supervisor.ToEntityStateShort(NewStates()),
	}
}
