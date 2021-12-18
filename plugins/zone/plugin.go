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
	"sync"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/plugins"
)

var (
	log = common.MustGetLogger("plugins.zone")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	actors     map[string]entity_manager.PluginActor
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]entity_manager.PluginActor),
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	return nil
}

// Unload ...
func (p plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
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

	actor := NewActor(entity, p.ScriptService, p.Adaptors, p.EventBus)
	p.actors[entity.Id.Name()] = actor
	p.EntityManager.Spawn(actor.Spawn)

	return
}

// RemoveActor ...
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

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginBuiltIn
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
		ActorCustomAttrs: false,
		ActorAttrs:       NewAttr(),
		ActorSetts:       NewSettings(),
	}
}
