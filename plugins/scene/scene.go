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

package scene

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugin_manager"
	"github.com/e154/smart-home/system/scripts"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.scene")
)

type pluginScene struct {
	entityManager *entity_manager.EntityManager
	eventBus      *event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]*EntityActor
	adaptors      *adaptors.Adaptors
	scriptService *scripts.ScriptService
}

func Register(manager *plugin_manager.PluginManager,
	entityManager *entity_manager.EntityManager,
	bus *event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	scriptService *scripts.ScriptService) {
	manager.Register(&pluginScene{
		entityManager: entityManager,
		eventBus:      bus,
		actorsLock:    &sync.Mutex{},
		actors:        make(map[string]*EntityActor),
		adaptors:      adaptors,
		scriptService: scriptService,
	})
	return
}

func (p *pluginScene) Load(service plugin_manager.IPluginManager, plugins map[string]interface{}) (err error) {

	if err := p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler); err != nil {
		log.Error(err.Error())
	}

	return
}

func (p *pluginScene) Unload() (err error) {

	return
}

func (p pluginScene) Name() string {
	return Name
}

func (p *pluginScene) AddOrUpdateEntity(entity *m.Entity) error {
	return p.addOrUpdateEntity(entity, entity.Attributes.Serialize())
}

func (p *pluginScene) RemoveEntity(entity *m.Entity) (err error) {
	return p.removeEntity(entity.Id.Name())
}

func (p *pluginScene) addOrUpdateEntity(entity *m.Entity,
	attributes m.EntityAttributeValue,
) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	name := entity.Id.Name()
	if _, ok := p.actors[name]; ok {
		return
	}

	if actor, ok := p.actors[name]; ok {
		// update
		actor.SetState(entity_manager.EntityStateParams{
			AttributeValues: attributes,
		})
		return
	}

	log.Infof("Add scene '%s'", name)

	var actor *EntityActor
	if actor, err = NewEntityActor(entity, attributes, p.adaptors, p.scriptService); err != nil {
		return
	}
	p.actors[name] = actor
	p.entityManager.Spawn(p.actors[name].Spawn)

	return
}

func (p *pluginScene) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		return
	}

	delete(p.actors, name)

	return
}

func (p *pluginScene) eventHandler(msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventCallScene:
		actor, ok := p.actors[v.EntityId.Name()]
		if !ok {
			return
		}
		actor.addEvent(v)

	default:
		//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
	}
}

func (p *pluginScene) Type() plugin_manager.PlugableType {
	return plugin_manager.PlugableBuiltIn
}

func (p *pluginScene) Depends() []string {
	return nil
}
