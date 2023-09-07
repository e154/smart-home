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

package zone

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.zone")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actorsLock *sync.Mutex
	actors     map[string]supervisor.PluginActor
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]supervisor.PluginActor),
	}
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	return nil
}

// Unload ...
func (p plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {

	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	attributes := entity.Attributes.Serialize()
	if actor, ok := p.actors[entity.Id.Name()]; ok {
		// update
		_ = actor.SetState(supervisor.EntityStateParams{
			AttributeValues: attributes,
		})
		return
	}

	actor := NewActor(entity, p.ScriptService, p.Adaptors, p.EventBus)
	p.actors[entity.Id.Name()] = actor
	p.Supervisor.Spawn(actor.Spawn)

	return
}

// RemoveActor ...
func (p *plugin) RemoveActor(entityId common.EntityId) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entityId.Name()]; !ok {
		err = errors.Wrap(apperr.ErrNotFound, fmt.Sprintf("failed remove '%s", entityId))
		return
	}

	delete(p.actors, entityId.Name())

	return
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
	return m.PluginOptions{
		ActorCustomAttrs: false,
		ActorAttrs:       NewAttr(),
		ActorSetts:       NewSettings(),
	}
}
