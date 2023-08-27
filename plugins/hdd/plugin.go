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

package hdd

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
	"sync"
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actorsLock *sync.Mutex
	actors     map[common.EntityId]*Actor
}

// New ...
func New() supervisor.Pluggable {
	p := &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*Actor),
	}
	return p
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}
	_ = p.EventBus.Subscribe("system/entities/+", p.eventHandler)
	return nil
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}
	_ = p.EventBus.Unsubscribe("system/entities/+", p.eventHandler)
	return nil
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		return
	}

	actor := NewActor(entity, p.Supervisor, p.EventBus)
	p.actors[entity.Id] = actor
	p.Supervisor.Spawn(actor.Spawn)

	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	delete(p.actors, entityId)

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
		actor, ok := p.actors[v.EntityId]
		if !ok {
			return
		}
		actor.runAction(v)
	}
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
		ActorAttrs:         NewAttr(),
		ActorSetts:         NewSettings(),
		ActorActions:       supervisor.ToEntityActionShort(NewActions()),
		ActorCustomActions: false,
		ActorCustomSetts:   false,
	}
}
