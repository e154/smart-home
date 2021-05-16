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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/plugins"
	"github.com/prometheus/common/log"
	"go.uber.org/atomic"
	"time"
)

var _ plugins.Plugable = (*plugin)(nil)

func init() {
	plugins.RegisterPlugin(Name, New)
}

type plugin struct {
	entityManager entity_manager.EntityManager
	isStarted     *atomic.Bool
	quit          chan struct{}
	pause         uint
	actor         *EntityActor
	adaptors      *adaptors.Adaptors
}

func New() plugins.Plugable {
	p := &plugin{
		isStarted: atomic.NewBool(false),
		pause:     50,
	}
	return p
}

func (c *plugin) Load(service plugins.Service) error {
	c.adaptors = service.Adaptors()
	c.entityManager = service.EntityManager()
	c.actor = NewEntityActor(c.entityManager)

	if c.isStarted.Load() {
		return nil
	}

	c.entityManager.Spawn(c.actor.Spawn)

	list, _, err := c.adaptors.Metric.Search("cpuspeed", 1, 0)
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
		if metric.Id, err = c.adaptors.Metric.Add(metric); err == nil {
			c.adaptors.Entity.AppendMetric(c.actor.Id, metric)
		}

	} else {
		metric = list[0]
	}

	c.actor.Metric = []m.Metric{metric}

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(c.pause))
		c.quit = make(chan struct{})
		defer func() {
			ticker.Stop()
			close(c.quit)
		}()

		c.isStarted.Store(true)
		defer func() {
			c.isStarted.Store(false)
		}()

		for {
			select {
			case <-c.quit:
				return
			case <-ticker.C:
				c.actor.selfUpdate()
			}
		}
	}()

	return nil
}

func (c *plugin) Unload() error {
	if !c.isStarted.Load() {
		return nil
	}
	c.quit <- struct{}{}
	return nil
}

func (c plugin) Name() string {
	return Name
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
