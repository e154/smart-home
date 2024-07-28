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

package memory

import (
	"sync"

	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	total       metrics.Gauge
	free        metrics.Gauge
	usedPercent metrics.GaugeFloat64
	updateLock  *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, service),
		total:       metrics.NewGauge(),
		free:        metrics.NewGauge(),
		usedPercent: metrics.NewGaugeFloat64(),
		updateLock:  &sync.Mutex{},
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

	v, _ := mem.VirtualMemory()
	e.total.Update(int64(v.Total))
	e.free.Update(int64(v.Free))

	usedPercent := common.Rounding32(v.UsedPercent, 2)
	e.usedPercent.Update(float64(usedPercent))

	e.AttrMu.Lock()
	e.Attrs[AttrTotal].Value = v.Total
	e.Attrs[AttrFree].Value = v.Free
	e.Attrs[AttrUsedPercent].Value = usedPercent
	e.AttrMu.Unlock()

	//e.SetMetric(e.Id, "memory", map[string]float32{
	//	"used_percent": usedPercent,
	//})

	e.SaveState(false, false)
}
