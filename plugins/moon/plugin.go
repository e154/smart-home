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
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

var (
	log = logger.MustGetLogger("plugins.moon")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actorsLock *sync.Mutex
	actors     map[string]*Actor
	quit       chan struct{}
	pause      time.Duration
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin:     supervisor.NewPlugin(),
		actorsLock: &sync.Mutex{},
		actors:     make(map[string]*Actor),
		pause:      240,
	}
}

// Load ...
func (p *plugin) Load(service supervisor.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.quit = make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second * p.pause)

		defer func() {
			ticker.Stop()
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

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}

	p.quit <- struct{}{}
	return nil
}

// Name ...
func (p *plugin) Name() string {
	return Name
}

// AddOrUpdateActor ...
func (p *plugin) AddOrUpdateActor(entity *m.Entity) (err error) {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	if _, ok := p.actors[entity.Id.Name()]; ok {
		p.actors[entity.Id.Name()].setPosition(entity.Settings)
		p.actors[entity.Id.Name()].UpdateMoonPosition(time.Now())
		return
	}

	p.actors[entity.Id.Name()] = NewActor(entity, p.Supervisor, p.Adaptors, p.ScriptService, p.EventBus)
	p.Supervisor.Spawn(p.actors[entity.Id.Name()].Spawn)

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

func (p *plugin) updatePositionForAll() {
	p.actorsLock.Lock()
	defer p.actorsLock.Unlock()

	now := time.Now()
	for _, actor := range p.actors {
		actor.UpdateMoonPosition(now)
	}
}

// Type ...
func (p *plugin) Type() supervisor.PluginType {
	return supervisor.PluginInstallable
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
		Actors:      true,
		ActorAttrs:  NewAttr(),
		ActorSetts:  NewSettings(),
		ActorStates: supervisor.ToEntityStateShort(NewStates()),
	}
}
