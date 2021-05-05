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

package modbus

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var (
	log = common.MustGetLogger("plugins.modbus")
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
	pause         time.Duration
	quit          chan struct{}
}

func New() plugins.Plugable {
	return &plugin{
		isStarted:  atomic.NewBool(false),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*EntityActor),
		pause:      240,
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
	case event_bus.EventStateChanged:
		if v.Type != "node" {
			return
		}

	case event_bus.EventRemoveEntity:
	}

	return
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
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
	//fmt.Println("updatePositionForAll")

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

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
