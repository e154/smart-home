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
	entityManager entity_manager.EntityManager
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

func (u *plugin) Load(service plugins.Service) error {
	u.adaptors = service.Adaptors()
	u.entityManager = service.EntityManager()

	if u.isStarted.Load() {
		return nil
	}
	u.isStarted.Store(true)

	u.entity = NewActor(u.entityManager)
	u.quit = make(chan struct{})

	u.storyModel = &m.RunStory{
		Start: time.Now(),
	}

	var err error
	u.storyModel.Id, err = u.adaptors.RunHistory.Add(u.storyModel)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	u.entityManager.Spawn(u.entity.Spawn)

	go func() {
		ticker := time.NewTicker(time.Second * u.pause)
		defer func() {
			ticker.Stop()
			u.isStarted.Store(false)
			close(u.quit)
		}()

		for {
			select {
			case <-u.quit:
				return
			case <-ticker.C:
				u.entity.update()
			}
		}
	}()
	return nil
}

func (u *plugin) Unload() (err error) {
	if !u.isStarted.Load() {
		return
	}
	u.quit <- struct{}{}
	u.storyModel.End = common.Time(time.Now())
	if err = u.adaptors.RunHistory.Update(u.storyModel); err != nil {
		log.Error(err.Error())
	}
	return
}

func (u plugin) Name() string {
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
