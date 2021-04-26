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

package updater

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/plugins"
	atomic2 "go.uber.org/atomic"
	"time"
)

const (
	name = "updater"
	uri  = "https://api.github.com/repos/e154/smart-home/releases/latest"
)

var (
	log = common.MustGetLogger("plugins.updater")
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	entityManager entity_manager.EntityManager
	isStarted     atomic2.Bool
	pause         time.Duration
	entity        *EntityActor
	quit          chan struct{}
}

func New() plugins.Plugable {
	return &plugin{
		pause: 24,
	}
}

func (u *plugin) Load(service plugins.Service) (err error) {
	u.entityManager = service.EntityManager()

	if u.isStarted.Load() {
		return
	}

	u.entity = NewEntityActor(u.entityManager)

	u.entityManager.Spawn(u.entity.Spawn)
	u.entity.check()
	u.quit = make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Hour * u.pause)

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
				u.entity.check()
			}
		}
	}()

	return
}

func (u *plugin) Unload() (err error) {
	if !u.isStarted.Load() {
		return
	}
	u.quit <- struct{}{}
	return
}

func (u *plugin) Name() string {
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
