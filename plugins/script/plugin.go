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

package script

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.script")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	plugins.Plugin
	entityManager entity_manager.EntityManager
	eventBus      event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]*Actor
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
}

func New() plugins.Plugable {
	return &plugin{
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	if err := p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler); err != nil {
		log.Error(err.Error())
	}

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)

	return
}

func (p plugin) Name() string {
	return Name
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

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	name := entity.Id.Name()
	if _, ok := p.actors[name]; ok {
		return
	}

	if actor, ok := p.actors[name]; ok {
		// update
		actor.SetState(entity_manager.EntityStateParams{
			AttributeValues: entity.Attributes.Serialize(),
			SettingsValue:   entity.Settings.Serialize(),
		})
		return
	}

	var actor *Actor
	if actor, err = NewActor(entity, p.adaptors, p.scriptService, p.entityManager, p.eventBus); err != nil {
		return
	}
	p.actors[name] = actor
	p.entityManager.Spawn(p.actors[name].Spawn)

	return
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventCallAction:
		actor, ok := p.actors[v.EntityId.Name()]
		if !ok {
			return
		}
		actor.addAction(v)

	default:
		//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
	}
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
		Actors:             true,
		ActorCustomAttrs:   true,
		ActorCustomStates:  true,
		ActorCustomActions: true,
	}
}
