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

package modbus_rtu

import (
	"fmt"
	"github.com/e154/smart-home/system/supervisor"
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
)

var (
	log = logger.MustGetLogger("plugins.modbus_rtu")
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
	return &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*Actor),
	}
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

// Name ...
func (p *plugin) Name() string {
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
		actor.addAction(v)
	}
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		return
	}

	actor := NewActor(entity, p.Supervisor, p.Adaptors, p.ScriptService, p.EventBus)
	p.actors[entity.Id] = actor
	p.Supervisor.Spawn(actor.Spawn)
	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", entityId))
		return
	}

	delete(p.actors, entityId)
	return
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
}

// Depends ...
func (p *plugin) Depends() []string {
	return []string{"node"}
}

// Version ...
func (p *plugin) Version() string {
	return "0.0.1"
}

// Options ...
func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		Actors:             true,
		ActorCustomAttrs:   true,
		ActorCustomActions: true,
		ActorCustomStates:  true,
		ActorAttrs:         NewAttr(),
		ActorSetts:         NewSettings(),
	}
}
