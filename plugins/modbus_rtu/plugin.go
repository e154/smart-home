// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/system/scripts"
	"go.uber.org/atomic"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.modbus_rtu")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	isStarted     *atomic.Bool
	eventBus      event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[common.EntityId]*EntityActor
}

func New() plugins.Plugable {
	return &plugin{
		isStarted:  atomic.NewBool(false),
		actorsLock: &sync.Mutex{},
		actors:     make(map[common.EntityId]*EntityActor),
	}
}

func (p *plugin) Load(service plugins.Service) error {
	p.adaptors = service.Adaptors()
	p.eventBus = service.EventBus()
	p.entityManager = service.EntityManager()
	p.scriptService = service.ScriptService()

	if p.isStarted.Load() {
		return nil
	}
	p.isStarted.Store(true)
	p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)

	return nil
}

func (p *plugin) Unload() error {
	if !p.isStarted.Load() {
		return nil
	}
	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)
	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventAddedNewEntity:
	case event_bus.EventStateChanged:
		if v.Type != "node" {
			return
		}

	case event_bus.EventRemoveEntity:
	}

	return
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id]; ok {
		err = fmt.Errorf("the actor with id '%s' has already been created", entity.Id)
		return
	}

	var actor *EntityActor
	actor = NewEntityActor(entity, p.entityManager, p.adaptors, p.scriptService, p.eventBus)
	p.actors[entity.Id] = actor
	p.entityManager.Spawn(actor.Spawn)
	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) error {
	return nil
}

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return []string{"node"}
}

func (p *plugin) Version() string {
	return "0.0.1"
}
