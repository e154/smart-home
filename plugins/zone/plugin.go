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

package zone

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugin_manager"
	"go.uber.org/atomic"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.zone")
)

type plugin struct {
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	eventBus      event_bus.EventBus
	isStarted     *atomic.Bool
	actorsLock    *sync.Mutex
	actors        map[string]entity_manager.PluginActor
}

func Register(manager plugin_manager.PluginManager,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors) {
	manager.Register(&plugin{
		entityManager: entityManager,
		isStarted:     atomic.NewBool(false),
		eventBus:      eventBus,
		adaptors:      adaptors,
		actorsLock:    &sync.Mutex{},
		actors:        make(map[string]entity_manager.PluginActor),
	})
}

func (p *plugin) Load(service plugin_manager.PluginManager, plugins map[string]interface{}) (err error) {
	return
}

func (p plugin) Unload() error {
	return nil
}

func (p plugin) Name() string {
	return Name
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	attributes := entity.Attributes.Serialize()
	if actor, ok := p.actors[entity.Id.Name()]; ok {
		// update
		actor.SetState(entity_manager.EntityStateParams{
			AttributeValues: attributes,
		})
		return
	}

	actor := NewEntityActor(entity.Id.Name(), attributes)
	p.actors[entity.Id.Name()] = actor
	p.entityManager.Spawn(actor.Spawn)

	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId.Name()]; !ok {
		err = fmt.Errorf("not found")
		return
	}

	delete(p.actors, entityId.Name())

	return
}

func (p *plugin) Type() plugin_manager.PlugableType {
	return plugin_manager.PlugableBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}
