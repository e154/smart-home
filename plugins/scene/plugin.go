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

package scene

import (
	"fmt"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"sync"
)

var (
	log = common.MustGetLogger("plugins.scene")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	actorsLock *sync.Mutex
	actors     map[string]*Actor
}

func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
	}
}

func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	if err := p.EventBus.Subscribe(event_bus.TopicEntities, p.eventHandler); err != nil {
		log.Error(err.Error())
	}

	return nil
}

func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.EventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)

	return
}

func (p plugin) Name() string {
	return Name
}

func (p *plugin) AddOrUpdateActor(entity *m.Entity) error {
	return p.addOrUpdateEntity(entity, entity.Attributes.Serialize())
}

func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	return p.removeEntity(entityId.Name())
}

func (p *plugin) addOrUpdateEntity(entity *m.Entity,
	attributes m.AttributeValue,
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

	var actor *Actor
	if actor, err = NewActor(entity, attributes,
		p.Adaptors, p.ScriptService, p.EntityManager); err != nil {
		return
	}
	p.actors[name] = actor
	p.EntityManager.Spawn(p.actors[name].Spawn)

	return
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

func (p *plugin) eventHandler(_ string, msg interface{}) {

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
	return m.PluginOptions{}
}
