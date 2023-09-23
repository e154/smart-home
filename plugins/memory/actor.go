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
	"sync"

	m "github.com/e154/smart-home/models"

	"github.com/e154/smart-home/common/events"

	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	supervisor.BaseActor
	total       metrics.Gauge
	free        metrics.Gauge
	usedPercent metrics.GaugeFloat64
	updateLock  *sync.Mutex
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor: supervisor.NewBaseActor(entity, service),
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

func (u *Actor) Destroy() {

}

// Spawn ...
func (u *Actor) Spawn() {
	return
}

func (u *Actor) selfUpdate() {

	u.updateLock.Lock()
	defer u.updateLock.Unlock()

	oldState := u.GetEventState()
	u.Now(oldState)

	v, _ := mem.VirtualMemory()
	u.total.Update(int64(v.Total))
	u.free.Update(int64(v.Free))

	usedPercent := common.Rounding32(v.UsedPercent, 2)
	u.usedPercent.Update(float64(usedPercent))

	u.AttrMu.Lock()
	u.Attrs[AttrTotal].Value = v.Total
	u.Attrs[AttrFree].Value = v.Free
	u.Attrs[AttrUsedPercent].Value = usedPercent
	u.AttrMu.Unlock()

	//u.SetMetric(u.Id, "memory", map[string]float32{
	//	"used_percent": usedPercent,
	//})

	u.Service.EventBus().Publish("system/entities/"+u.Id.String(), events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(),
	})
}
