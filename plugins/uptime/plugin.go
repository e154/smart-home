// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package uptime

import (
	"context"
	"fmt"
	"time"

	"github.com/e154/smart-home/system/supervisor"

	"github.com/e154/smart-home/common/logger"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	name = "uptime"
)

var (
	log = logger.MustGetLogger("plugins.uptime")
)

var _ supervisor.Pluggable = (*plugin)(nil)

func init() {
	supervisor.RegisterPlugin(Name, New)
}

type plugin struct {
	*supervisor.Plugin
	actor      *Actor
	ticker     *time.Ticker
	storyModel *m.RunStory
}

// New ...
func New() supervisor.Pluggable {
	return &plugin{
		Plugin: supervisor.NewPlugin(),
	}
}

// Load ...
func (p *plugin) Load(ctx context.Context, service supervisor.Service) (err error) {
	if err = p.Plugin.Load(ctx, service, nil); err != nil {
		return
	}

	var entity = &m.Entity{
		Id:         common.EntityId(fmt.Sprintf("%s.%s", EntitySensor, Name)),
		PluginName: Name,
		Attributes: NewAttr(),
	}
	p.actor = NewActor(entity, service)

	p.storyModel = &m.RunStory{
		Start: time.Now(),
	}
	if err = p.AddActor(p.actor, entity); err != nil {
		return
	}

	p.storyModel.Id, err = p.Service.Adaptors().RunHistory.Add(context.Background(), p.storyModel)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	go func() {
		const pause = 60
		p.ticker = time.NewTicker(time.Second * pause)

		for range p.ticker.C {
			p.actor.update()
		}
	}()
	return nil
}

// Unload ...
func (p *plugin) Unload(ctx context.Context) (err error) {
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}
	if err = p.Plugin.Unload(ctx); err != nil {
		return
	}

	p.storyModel.End = common.Time(time.Now())
	if err = p.Service.Adaptors().RunHistory.Update(context.Background(), p.storyModel); err != nil {
		log.Error(err.Error())
	}
	return
}

// Name ...
func (p plugin) Name() string {
	return name
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
