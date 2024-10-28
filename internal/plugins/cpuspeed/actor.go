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

package cpuspeed

import (
	"sync"

	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/v4/cpu"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	cores           int64
	model           string
	mhz             float64
	all             metrics.GaugeFloat64
	loadMin         metrics.GaugeFloat64
	loadMax         metrics.GaugeFloat64
	allCpuPrevTotal float64
	allCpuPrevIdle  float64
	updateLock      *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) *Actor {

	actor := &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		all:        metrics.NewGaugeFloat64(),
		loadMin:    metrics.NewGaugeFloat64(),
		loadMax:    metrics.NewGaugeFloat64(),
		updateLock: &sync.Mutex{},
	}

	cpuInfo, err := cpu.Info()
	if err == nil {
		actor.mhz = cpuInfo[0].Mhz
		actor.cores = int64(cpuInfo[0].Cores)
		actor.model = cpuInfo[0].Model
	}

	if entity != nil {
		actor.Metric = entity.Metrics
	}

	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) selfUpdate() {

	e.updateLock.Lock()
	defer e.updateLock.Unlock()

	// export CGO_ENABLED=1
	timeStats, err := cpu.Times(false)
	if err != nil {
		return
	}
	if len(timeStats) == 0 {
		return
	}

	total := timeStats[0].Total()
	diffIdle := timeStats[0].Idle - e.allCpuPrevIdle
	diffTotal := total - e.allCpuPrevTotal
	all := common.Rounding(100*(diffTotal-diffIdle)/diffTotal, 2)
	if v := e.loadMin.Value(); v == 0 || all < v {
		e.loadMin.Update(all)
	}
	if v := e.loadMax.Value(); v == 0 || all > v {
		e.loadMax.Update(all)
	}

	e.all.Update(all)
	e.allCpuPrevTotal = total
	e.allCpuPrevIdle = timeStats[0].Idle

	e.AttrMu.Lock()
	e.Attrs[AttrCpuCores].Value = e.cores
	e.Attrs[AttrCpuMhz].Value = e.mhz
	e.Attrs[AttrCpuAll].Value = e.all.Value()
	e.Attrs[AttrLoadMax].Value = e.loadMax.Value()
	e.Attrs[AttrLoadMin].Value = e.loadMin.Value()
	e.AttrMu.Unlock()

	//e.SetMetric(e.Id, "cpuspeed", map[string]float32{
	//	"all": common.Rounding32(e.all.Value(), 2),
	//})

	e.SaveState(false, false)
}
