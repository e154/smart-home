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

// +build linux,!mips64,!mips64le darwin

package uptime

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/plugins"
	"go.uber.org/atomic"
	"time"
)

const (
	name = "uptime"
)

var (
	log = common.MustGetLogger("plugins.uptime")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	plugins.Plugin
	entityManager entity_manager.EntityManager
	eventBus          event_bus.EventBus
	entity        *Actor
	isStarted     *atomic.Bool
	ticker        *time.Ticker
	pause         time.Duration
	adaptors      *adaptors.Adaptors
	storyModel    *m.RunStory
	quit          chan struct{}
}

func New() plugins.Plugable {
	return &plugin{
		isStarted: atomic.NewBool(false),
		pause:     60,
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

	p.entity = NewActor(p.entityManager, p.eventBus)
	p.quit = make(chan struct{})

	p.storyModel = &m.RunStory{
		Start: time.Now(),
	}

	var err error
	p.storyModel.Id, err = p.adaptors.RunHistory.Add(p.storyModel)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	p.entityManager.Spawn(p.entity.Spawn)

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
				p.entity.update()
			}
		}
	}()
	return nil
}

func (p *plugin) Unload() (err error) {
	if !p.isStarted.Load() {
		return
	}
	p.quit <- struct{}{}
	p.storyModel.End = common.Time(time.Now())
	if err = p.adaptors.RunHistory.Update(p.storyModel); err != nil {
		log.Error(err.Error())
	}
	return
}

func (p plugin) Name() string {
	return name
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
