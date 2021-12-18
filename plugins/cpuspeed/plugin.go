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

package cpuspeed

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
	quit  chan struct{}
	pause uint
	actor *Actor
}

// New ...
func New() plugins.Plugable {
	p := &plugin{
		Plugin: plugins.NewPlugin(),
		pause:  50,
	}
	return p
}

// Load ...
func (p *plugin) Load(service plugins.Service) (err error) {
	if err = p.Plugin.Load(service); err != nil {
		return
	}

	p.EntityManager.Spawn(p.actor.Spawn)

	list, _, err := p.Adaptors.Metric.Search("cpuspeed", 1, 0)
	if err != nil {
		log.Error(err.Error())
	}

	var metric m.Metric
	if len(list) == 0 {
		metric = m.Metric{
			Name:        "cpuspeed",
			Description: "Cpu metric",
			Options: m.MetricOptions{
				Items: []m.MetricOptionsItem{
					{
						Name:        "mhz",
						Description: "",
						Color:       "#ff0000",
						Translate:   "mhz",
						Label:       "GHz",
					},
					{
						Name:        "all",
						Description: "",
						Color:       "#00ff00",
						Translate:   "all",
						Label:       "GHz",
					},
				},
			},
			Type: common.MetricTypeLine,
		}
		if metric.Id, err = p.Adaptors.Metric.Add(metric); err == nil {
			p.Adaptors.Entity.AppendMetric(p.actor.Id, metric)
		}

	} else {
		metric = list[0]
	}

	p.actor.Metric = []m.Metric{metric}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(p.pause))
		p.quit = make(chan struct{})
		defer func() {
			ticker.Stop()
			close(p.quit)
		}()

		for {
			select {
			case <-p.quit:
				return
			case <-ticker.C:
				p.actor.selfUpdate()
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
func (p plugin) Name() string {
	return Name
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
		ActorAttrs: NewAttr(),
	}
}
