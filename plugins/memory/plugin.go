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

package memory

import (
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/plugins"
	"github.com/prometheus/common/log"
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	*plugins.Plugin
	ticker *time.Ticker
	pause  uint
	actor  *Actor
}

// New ...
func New() plugins.Plugable {
	p := &plugin{
		Plugin: plugins.NewPlugin(),
		pause:  10,
	}
	return p
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.actor = NewActor(service.EntityManager(), service.EventBus())
	p.EntityManager.Spawn(p.actor.Spawn)

	list, _, err := p.Adaptors.Metric.Search("memory", 1, 0)
	if err != nil {
		log.Error(err.Error())
	}

	var metric *m.Metric
	if len(list) == 0 {
		metric = &m.Metric{
			Name:        "memory",
			Description: "RAM metric",
			Options: m.MetricOptions{
				Items: []m.MetricOptionsItem{
					{
						Name:        "used_percent",
						Description: "",
						Color:       "#0000FF",
						Translate:   "used_percent",
						Label:       "%",
					},
				},
			},
			Type: common.MetricTypeLine,
		}
		if metric.Id, err = p.Adaptors.Metric.Add(metric); err == nil {
			_ = p.Adaptors.Entity.AppendMetric(p.actor.Id, metric)
		}

	} else {
		metric = list[0]
	}

	p.actor.Metric = []*m.Metric{metric}

	go func() {
		p.ticker = time.NewTicker(time.Second * time.Duration(p.pause))

		for range p.ticker.C {
			p.actor.selfUpdate()
		}
	}()

	return nil
}

// Unload ...
func (p *plugin) Unload() (err error) {
	if err = p.Plugin.Unload(); err != nil {
		return
	}
	if p.ticker != nil {
		p.ticker.Stop()
		p.ticker = nil
	}
	return nil
}

// Name ...
func (p plugin) Name() string {
	return Name
}

// Type ...
func (p *plugin) Type() plugins.PluginType {
	return plugins.PluginInstallable
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
		ActorAttrs: NewAttr(),
	}
}
