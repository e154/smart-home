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
	"context"
	"fmt"
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.scene")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actorsLock *sync.Mutex
	actors     map[string]*Actor
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
	}
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	if err := p.EventBus.Subscribe("system/entities/+", p.eventHandler); err != nil {
		log.Error(err.Error())
	}

	return nil
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	_ = p.EventBus.Unsubscribe("system/entities/+", p.eventHandler)

	return
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) error {
	return p.addOrUpdateEntity(entity, entity.Attributes.Serialize())
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	return p.removeEntity(entityId.Name())
}

func (p *plugin) addOrUpdateEntity(entity *m.Entity, attributes m.AttributeValue) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	name := entity.Id.Name()
	if _, ok := p.actors[name]; ok {
		return
	}

	if len(entity.Actions) == 0 {
		var action = &m.EntityAction{
			Name:        "apply",
			Description: "apply scene",
			EntityId:    entity.Id,
		}
		if action.Id, err = p.Adaptors.EntityAction.Add(context.Background(), action); err != nil {
			return
		}
	}

	if actor, ok := p.actors[name]; ok {
		// update
		_ = actor.SetState(supervisor.EntityStateParams{
			AttributeValues: attributes,
		})
		return
	}

	var actor *Actor
	if actor, err = NewActor(entity, attributes,
		p.Adaptors, p.ScriptService, p.Supervisor); err != nil {
		return
	}
	p.actors[name] = actor
	p.Supervisor.Spawn(p.actors[name].Spawn)

	return
}

func (p *plugin) removeEntity(name string) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[name]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove \"%s\"", name))
		return
	}

	delete(p.actors, name)

	return
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventCallEntityAction:
		actor, ok := p.actors[v.EntityId.Name()]
		if !ok {
			return
		}
		actor.addEvent(events.EventCallScene{
			PluginName: v.PluginName,
			EntityId:   v.EntityId,
			Args:       v.Args,
		})
	case events.EventCallScene:
		actor, ok := p.actors[v.EntityId.Name()]
		if !ok {
			return
		}
		actor.addEvent(v)

	default:
		//fmt.Printf("new event: %v\n", reflect.TypeOf(v).String())
	}
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginBuiltIn
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
	return m.PluginOptions{}
}
