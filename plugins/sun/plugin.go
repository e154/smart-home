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

package sun

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/zone"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugin_manager"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var (
	log = common.MustGetLogger("plugins.sun")
)

type plugin struct {
	entityManager entity_manager.EntityManager
	adaptors      *adaptors.Adaptors
	isStarted     *atomic.Bool
	eventBus      event_bus.EventBus
	actorsLock    *sync.Mutex
	actors        map[string]*EntityActor
	pause         time.Duration
	quit          chan struct{}
}

func Register(manager plugin_manager.PluginManager,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus,
	adaptors *adaptors.Adaptors,
	second time.Duration) {
	manager.Register(&plugin{
		entityManager: entityManager,
		adaptors:      adaptors,
		isStarted:     atomic.NewBool(false),
		eventBus:      eventBus,
		actorsLock:    &sync.Mutex{},
		actors:        make(map[string]*EntityActor),
		pause:         second,
	})
}

func (p *plugin) Load(service plugin_manager.PluginManager, plugins map[string]interface{}) error {
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

func (p *plugin) eventHandler(msg interface{}) {

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

		entityId := common.EntityId(fmt.Sprintf("sun.%s", v.EntityId.Name()))
		p.entityManager.Remove(entityId)
	}

	return
}

func (p *plugin) addOrUpdateEntity(zoneName string, zoneAttr m.EntityAttributes) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	var lat, lon, elevation float64
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
	}

	if lat == 0 && lon == 0 {
		return
	}

	if _, ok := p.actors[zoneName]; ok {
		p.actors[zoneName].setPosition(lat, lon, elevation)
		p.actors[zoneName].updateSunPosition()
		return
	}

	p.actors[zoneName] = NewEntityActor(zoneName)
	p.entityManager.Spawn(p.actors[zoneName].Spawn)

	if zoneAttr != nil {
		p.actors[zoneName].setPosition(lat, lon, elevation)
		p.actors[zoneName].updateSunPosition()
	}

	return
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	return p.addOrUpdateEntity(entity.Id.Name(), nil)
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
	//fmt.Println("updatePositionForAll")

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	for _, actor := range p.actors {
		actor.updateSunPosition()
	}
}

func (p *plugin) Type() plugin_manager.PlugableType {
	return plugin_manager.PlugableBuiltIn
}

func (p *plugin) Depends() []string {
	return nil
}
