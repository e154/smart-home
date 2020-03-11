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

package metric

import (
	"github.com/DrmagicE/gmqtt"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/logging"
	atomic2 "go.uber.org/atomic"
	"sync/atomic"
	"time"
)

var (
	log = common.MustGetLogger("metric")
)

const name = "metric"
const metricPrefix = "gmqtt_"

type Metric struct {
	metric       *metrics.MetricManager
	statsManager gmqtt.StatsManager
	path         string
	isStarted    atomic2.Bool
	quit         chan struct{}
	pause        int64
}

func New(metric *metrics.MetricManager, pause int64) *Metric {

	return &Metric{
		metric: metric,
		quit:   make(chan struct{}),
		pause:  pause,
	}
}

func (p *Metric) Load(service gmqtt.Server) (err error) {

	p.statsManager = service.GetStatsManager()

	if p.isStarted.Load() {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(p.pause))
		defer ticker.Stop()

		p.isStarted.Store(true)
		defer func() {
			p.isStarted.Store(false)
		}()

		for {
			select {
			case _, ok := <-p.quit:
				if !ok {
					return
				}
				return
			case <-ticker.C:
				p.collect()
			}
		}
	}()

	return
}

func (p *Metric) Unload() (err error) {
	if !p.isStarted.Load() {
		return
	}
	p.quit <- struct{}{}

	return
}

func (p *Metric) HookWrapper() gmqtt.HookWrapper {
	return gmqtt.HookWrapper{}
}

func (p *Metric) Name() string {
	return name
}

func (p *Metric) collect() {
	st := p.statsManager.GetStats()
	p.collectClientStats(st.ClientStats)
}

func (p *Metric) collectClientStats(c *gmqtt.ClientStats) {
	go p.metric.Update(metrics.MqttClientStats{
		ConnectedTotal:    atomic.LoadUint64(&c.ConnectedTotal),
		DisconnectedTotal: atomic.LoadUint64(&c.DisconnectedTotal),
		ActiveCurrent:     atomic.LoadUint64(&c.ActiveCurrent),
		InactiveCurrent:   atomic.LoadUint64(&c.InactiveCurrent),
		ExpiredTotal:      atomic.LoadUint64(&c.ExpiredTotal),
	})
}
