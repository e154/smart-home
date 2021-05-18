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

package moon

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
	"time"
)

var (
	log = common.MustGetLogger("plugins.moon")
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
	actors        map[string]*Actor
	quit          chan struct{}
	pause         time.Duration
}

func New() plugins.Plugable {
	return &plugin{
		isStarted:  atomic.NewBool(false),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
		pause:      240,
	}
}

func (p *plugin) Load(service plugins.Service) error {
	p.entityManager = service.EntityManager()
	p.eventBus = service.EventBus()
	p.adaptors = service.Adaptors()

	if p.isStarted.Load() {
		return nil
	}
	p.isStarted.Store(true)
	p.eventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)
	p.quit = make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second * p.pause)

		defer func() {
			ticker.Stop()
			p.isStarted.Store(false)
			close(p.quit)
		}()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				p.updatePositionForAll()
			}
		}
	}()

	return nil
}

func (p *plugin) Unload() error {
	if !p.isStarted.Load() {
		return nil
	}
	p.eventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)
	p.quit <- struct{}{}
	return nil
}

func (p *plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case event_bus.EventAddedNewEntity:
		if v.Type != "zone" {
			return
		}

		p.addOrUpdateEntity(v.EntityId.Name(), v.Attributes)

	case event_bus.EventStateChanged:
		if v.Type != "zone" {
			return
		}

		zoneAttr := zone.NewAttr()
		zoneAttr.Deserialize(v.NewState.Attributes.Serialize())
		p.addOrUpdateEntity(v.EntityId.Name(), zoneAttr)

	case event_bus.EventRemoveEntity:
		if v.Type != "zone" {
			return
		}

		if err := p.removeEntity(v.EntityId.Name()); err != nil {
			return
		}

		entityId := common.EntityId(fmt.Sprintf("moon.%s", v.EntityId.Name()))
		p.entityManager.Remove(entityId)
	}

	return
}

func (p *plugin) addOrUpdateEntity(name string, zoneAttr m.EntityAttributes) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	var lat, lon, elevation float64
	var timezone int
	if zoneAttr != nil {
		// lat
		if _lat, ok := zoneAttr[zone.AttrLat]; ok {
			lat, ok = _lat.Value.(float64)
		}

		// lon
		if _lon, ok := zoneAttr[zone.AttrLon]; ok {
			lon, ok = _lon.Value.(float64)
		}

		// elevation
		if _elevation, ok := zoneAttr[zone.AttrElevation]; ok {
			elevation, ok = _elevation.Value.(float64)
		}

		// timezone
		if _timezone, ok := zoneAttr[zone.AttrTimezone]; ok {
			timezone, ok = _timezone.Value.(int)
		}
	}

	if _, ok := p.actors[name]; ok {
		p.actors[name].setPosition(lat, lon, elevation, timezone)
		p.actors[name].updateMoonPosition()
		return
	}

	p.actors[name] = NewActor(name, p.entityManager)
	p.entityManager.Spawn(p.actors[name].Spawn)

	if zoneAttr != nil {
		p.actors[name].setPosition(lat, lon, elevation, timezone)
		p.actors[name].updateMoonPosition()
	}

	return
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.addOrUpdateEntity(entity.Id.Name(), entity.Attributes)
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

func (p *plugin) updatePositionForAll() {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		actor.updateMoonPosition()
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
