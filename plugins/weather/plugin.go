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
	"sync"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/event_bus/events"
	"github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/tmp/apperr"
)

const (
	// Name ...
	Name = "weather"
	// EntityWeather ...
	EntityWeather = string("weather")
)

var (
	log = logger.MustGetLogger("plugins.weather")
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

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin:     plugins.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	_ = p.EventBus.Subscribe(event_bus.TopicEntities, p.eventHandler)

	return nil
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	_ = p.EventBus.Unsubscribe(event_bus.TopicEntities, p.eventHandler)

	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

func (p *plugin) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventPassAttributes:
		if v.To.PluginName() != Name {
			return
		}

		_ = p.AddOrUpdateForecast(v.To.Name(), v.Attributes)
	}
}

// AddOrUpdateForecast ...
func (p *plugin) AddOrUpdateForecast(name string, attr m.Attributes) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	actor, ok := p.actors[name]
	if !ok {
		log.Warnf("forecast '%s.%s' not found", Name, name)
		return
	}

	var stateName string

	if a, ok := attr[AttrWeatherMain]; ok {
		stateName = a.String()
	}

	_ = actor.SetState(entity_manager.EntityStateParams{
		NewState:        common.String(stateName),
		AttributeValues: attr.Serialize(),
	})

	return
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	name := entity.Id.Name()
	if _, ok := p.actors[name]; !ok {
		p.actors[name] = NewActor(entity, p.EntityManager, p.EventBus)
		p.EntityManager.Spawn(p.actors[name].Spawn)
	}
	p.actors[name].UpdatePosition(entity.Settings)
	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) error {
	return p.removeEntity(entityId.Name())
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
		ActorAttrs:  BaseForecast(),
		ActorSetts:  NewSettings(),
		ActorStates: entity_manager.ToEntityStateShort(NewActorStates(false, false)),
	}
}
