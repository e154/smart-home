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
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/plugins"
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
	*plugins.Plugin
	entity     *Actor
	ticker     *time.Ticker
	pause      time.Duration
	storyModel *m.RunStory
	quit       chan struct{}
}

// New ...
func New() plugins.Plugable {
	return &plugin{
		Plugin: plugins.NewPlugin(),
		pause:  60,
	}
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.entity = NewActor(p.EntityManager, p.EventBus)
	p.quit = make(chan struct{})

	p.storyModel = &m.RunStory{
		Start: time.Now(),
	}

	p.storyModel.Id, err = p.Adaptors.RunHistory.Add(p.storyModel)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	p.EntityManager.Spawn(p.entity.Spawn)

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
				p.entity.update()
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
	p.storyModel.End = common.Time(time.Now())
	if err = p.Adaptors.RunHistory.Update(p.storyModel); err != nil {
		log.Error(err.Error())
	}
	return
}

// Name ...
func (p plugin) Name() string {
	return name
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
