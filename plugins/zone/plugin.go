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

package zone

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.zone")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	plugins.Plugin
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	eventBus      event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]entity_manager.PluginActor
}

func New() plugins.Plugable {
	return &plugin{
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]entity_manager.PluginActor),
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	return nil
}

func (p plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

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

	actor := NewActor(entity, p.ScriptService, p.Adaptors, p.eventBus)
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

func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}

func (p *plugin) Version() string {
	return "0.0.1"
}

func (p *plugin) Options() m.PluginOptions {
	return m.PluginOptions{
		ActorCustomAttrs: false,
		ActorAttrs:       NewAttr(),
		ActorSetts:       NewSettings(),
	}
}
