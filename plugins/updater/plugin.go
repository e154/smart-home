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

package updater

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	atomic2 "go.uber.org/atomic"
	"time"
)

const (
	name = "updater"
	uri  = "https://api.github.com/repos/e154/smart-home/releases/latest"
)

var (
	log = common.MustGetLogger("plugins.updater")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	plugins.Plugin
	entityManager entity_manager.EntityManager
	isStarted     atomic2.Bool
	pause         time.Duration
	actor         *Actor
	eventBus      event_bus.EventBus
	quit          chan struct{}
}

func New() plugins.Plugable {
	return &plugin{
		pause: 24,
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	p.entityManager = service.EntityManager()
	p.eventBus = service.EventBus()

	if p.isStarted.Load() {
		return
	}

	p.actor = NewActor(p.entityManager, p.eventBus)

	p.entityManager.Spawn(p.actor.Spawn)
	p.actor.check()
	p.quit = make(chan struct{})

	p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)

	go func() {
		ticker := time.NewTicker(time.Hour * p.pause)

		defer func() {
			ticker.Stop()
			p.isStarted.Store(false)
			close(p.quit)
		}()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				p.actor.check()
			}
		}
	}()

	return
}

func (p *plugin) Unload() (err error) {
	if !p.isStarted.Load() {
		return
	}
	p.quit <- struct{}{}
	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)
	return
}

func (p *plugin) Name() string {
	return name
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventCallAction:
		if v.EntityId != p.actor.Id {
			return
		}

		if v.ActionName == "check" {
			p.actor.check()
		}
	}

	return
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorAttrs:   NewAttr(),
		ActorActions: entity_manager.ToEntityActionShort(NewActions()),
		ActorStates:  entity_manager.ToEntityStateShort(NewStates()),
	}
}
