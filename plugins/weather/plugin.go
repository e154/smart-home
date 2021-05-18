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

package weather

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
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

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	isStarted     *atomic.Bool
	eventBus      event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]*EntityActor
}

func New() plugins.Plugable {
	return &plugin{
		isStarted:  atomic.NewBool(false),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*EntityActor),
	}
}

func (p *plugin) Load(service plugins.Service) error {
	p.adaptors = service.Adaptors()
	p.eventBus = service.EventBus()
	p.entityManager = service.EntityManager()

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

	p.isStarted.Store(false)

	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)

	return nil
}

func (p plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

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

		if err := p.removeEntity(v.EntityId.Name()); err != nil {
			return
		}

		entityId := common.EntityId(fmt.Sprintf("weather.%s", v.EntityId.Name()))
		p.entityManager.Remove(entityId)
	}

	return
}

func (p *plugin) addOrUpdateForecast(name string, attr m.EntityAttributes) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	actor, ok := p.actors[name]
	if !ok {
		log.Warnf("forecast '%s.%s' not found", Name, name)
		return
	}

	var stateName string

	if a, ok := attr[AttrWeatherCondition]; ok {
		stateName = a.String()
	}

	actor.SetState(entity_manager.EntityStateParams{
		NewState:        common.String(stateName),
		AttributeValues: attr.Serialize(),
	})

	return
}

func (p *plugin) addOrUpdateZone(name string, zoneAttr m.EntityAttributes) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		p.actors[name] = NewEntityActor(name, p.eventBus, p.entityManager)
		p.entityManager.Spawn(p.actors[name].Spawn)
	}
	p.actors[name].setPosition(zoneAttr)

	return
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.addOrUpdateZone(entity.Id.Name(), entity.Attributes)
	return
}

func (p *plugin) RemoveActor(entityId common.EntityId) error {
	return p.removeEntity(entityId.Name())
}

func (p *plugin) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		err = fmt.Errorf("not found")
		return
	}

	delete(p.actors, name)

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
