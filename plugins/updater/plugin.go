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

package updater

import (
	"context"
	"embed"
	"fmt"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

const (
	name = "updater"
	uri  = "https://api.github.com/repos/e154/smart-home/releases/latest"
)

var (
	log = logger.MustGetLogger("plugins.updater")
)

var _ supervisor.Pluggable = (*plugin)(nil)

//go:embed Readme.md
//go:embed Readme.ru.md
var F embed.FS

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	ticker *time.Ticker
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin: supervisor.NewPlugin(),
	}
	p.F = F
	return p
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, p.ActorConstructor); err != nil {
		return
	}

	if _, err = p.Service.Adaptors().Entity.GetById(context.Background(), common.EntityId(fmt.Sprintf("%s.%s", EntityUpdater, Name))); err != nil {
		entity := &m.Entity{
			Id:         common.EntityId(fmt.Sprintf("%s.%s", EntityUpdater, Name)),
			PluginName: Name,
			Attributes: NewAttr(),
		}
		err = p.Service.Adaptors().Entity.Add(context.Background(), entity)
	}

	_ = p.Service.EventBus().Subscribe("system/entities/+", p.eventHandler)

	go func() {
		const pause = 24
		p.ticker = time.NewTicker(time.Hour * pause)

		for range p.ticker.C {
			p.Actors.Range(func(key, value any) bool {
				actor, _ := value.(*Actor)
				actor.check()
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
	}

	_ = p.Plugin.Unload(ctx)

	_ = p.Service.EventBus().Unsubscribe("system/entities/+", p.eventHandler)
	return
}

// ActorConstructor ...
func (p *plugin) ActorConstructor(entity *m.Entity) (actor supervisor.PluginActor, err error) {
	actor = NewActor(entity, p.Service)
	return
}

// Name ...
func (p *plugin) Name() string {
	return name
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginBuiltIn
}

// Depends ...
func (p *plugin) Depends() []string {
	return nil
}

// Version ...
func (p *plugin) Version() string {
	return Version
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallEntityAction:
		values, ok := p.Check(v)
		if !ok {
			return
		}
		for _, value := range values {
			actor := value.(*Actor)
			actor.check()
		}
	}
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorAttrs:   NewAttr(),
		ActorActions: supervisor.ToEntityActionShort(NewActions()),
		ActorStates:  supervisor.ToEntityStateShort(NewStates()),
	}
}
