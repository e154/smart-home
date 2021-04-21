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

package weather

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugin_manager"
	"go.uber.org/atomic"
	"sync"
)

const (
	Name = "weather"
	// EntityWeather ...
	EntityWeather = common.EntityType("weather")
)

var (
	log = common.MustGetLogger("plugins.weather")
)

type pluginWeather struct {
	entityManager *entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	isStarted     *atomic.Bool
	eventBus      *event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]*EntityActor
}

func Register(manager *plugin_manager.PluginManager,
	entityManager *entity_manager.EntityManager,
	eventBus *event_bus.EventBus,
	adaptors *adaptors.Adaptors) {
	manager.Register(&pluginWeather{
		entityManager: entityManager,
		adaptors:      adaptors,
		isStarted:     atomic.NewBool(false),
		eventBus:      eventBus,
		actorsLock:    &sync.Mutex{},
		actors:        make(map[string]*EntityActor),
	})
}

func (p *pluginWeather) Load(service plugin_manager.IPluginManager, plugins map[string]interface{}) error {

	if p.isStarted.Load() {
		return nil
	}

	p.isStarted.Store(true)

	p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)

	return nil
}

func (p *pluginWeather) Unload() error {

	if !p.isStarted.Load() {
		return nil
	}

	p.isStarted.Store(false)

	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)

	return nil
}

func (p pluginWeather) Name() string {
	return Name
}

func (p *pluginWeather) eventHandler(msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventRequestState:
		if v.To.Type() != Name {
			return
		}

		p.addOrUpdateForecast(v.To.Name(), v.Attributes)

	case event_bus.EventAddedNewEntity:
		if v.Type != "zone" {
			return
		}

		p.addOrUpdateZone(v.EntityId.Name(), v.Attributes)

	case event_bus.EventStateChanged:
		if v.Type != "zone" {
			return
		}

		zoneAttr := zone.NewAttr()
		zoneAttr.Deserialize(v.NewState.Attributes.Serialize())
		p.addOrUpdateZone(v.EntityId.Name(), zoneAttr)

	case event_bus.EventRemoveEntity:
		if v.Type != "zone" {
			return
		}

		p.removeEntity(v.EntityId.Name())
	}

	return
}

func (p *pluginWeather) addOrUpdateForecast(name string, attr m.EntityAttributes) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	actor, ok := p.actors[name]
	if !ok {
		log.Warnf("forecast '%s.%s' not found", Name, name)
		return
	}

	actor.SetState(entity_manager.EntityStateParams{
		AttributeValues: attr.Serialize(),
	})

	return
}

func (p *pluginWeather) addOrUpdateZone(name string, zoneAttr m.EntityAttributes) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		log.Infof("Add weather '%v'", name)
	}

	if _, ok := p.actors[name]; !ok {
		p.actors[name] = NewEntityActor(name, p.eventBus)
		p.entityManager.Spawn(p.actors[name].Spawn)
	}
	p.actors[name].setPosition(zoneAttr)

	return
}

func (p *pluginWeather) AddOrUpdateEntity(entity *m.Entity) (err error) {
	p.addOrUpdateZone(entity.Id.Name(), entity.Attributes)
	return
}

func (p *pluginWeather) RemoveEntity(entity *m.Entity) error {
	return p.removeEntity(entity.Id.Name())
}

func (p *pluginWeather) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		return
	}

	delete(p.actors, name)

	return
}

func (p *pluginWeather) Type() plugin_manager.PlugableType {
	return plugin_manager.PlugableBuiltIn
}

func (p *pluginWeather) Depends() []string {
	return nil
}
