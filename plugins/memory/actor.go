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
	"fmt"
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	cores           int64
	model           string
	total           metrics.Gauge
	free            metrics.Gauge
	usedPercent     metrics.GaugeFloat64
	allCpuPrevTotal float64
	allCpuPrevIdle  float64
	eventBus        bus.Bus
	updateLock      *sync.Mutex
}

// NewActor ...
func NewActor(entityManager entity_manager.EntityManager,
	eventBus bus.Bus) *Actor {

	actor := &Actor{
		BaseActor: entity_manager.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityMemory, Name)),
			Name:              Name,
			EntityType:        EntityMemory,
			UnitOfMeasurement: "GHz",
			AttrMu:            &sync.RWMutex{},
			Attrs:             NewAttr(),
			Manager:           entityManager,
		},
		eventBus:    eventBus,
		total:       metrics.NewGauge(),
		free:        metrics.NewGauge(),
		usedPercent: metrics.NewGaugeFloat64(),
		updateLock:  &sync.Mutex{},
	}

	return actor
}

// Spawn ...
func (e *Actor) Spawn() entity_manager.PluginActor {
	return e
}

func (u *Actor) selfUpdate() {

	u.updateLock.Lock()
	defer u.updateLock.Unlock()

	oldState := u.GetEventState(u)
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

	u.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
