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
	"fmt"
	"sync"

	"github.com/e154/smart-home/common/events"

	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/v3/cpu"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/entity_manager"
)

// Actor ...
type Actor struct {
	entity_manager.BaseActor
	cores           int64
	model           string
	mhz             float64
	all             metrics.GaugeFloat64
	loadMin         metrics.GaugeFloat64
	loadMax         metrics.GaugeFloat64
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
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityCpuspeed, Name)),
			Name:              Name,
			EntityType:        EntityCpuspeed,
			UnitOfMeasurement: "GHz",
			AttrMu:            &sync.RWMutex{},
			Attrs:             NewAttr(),
			Manager:           entityManager,
		},
		eventBus:   eventBus,
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

	// export CGO_ENABLED=1
	timeStats, err := cpu.Times(false)
	if err != nil {
		return
	}
	if len(timeStats) == 0 {
		return
	}

	total := timeStats[0].Total()
	diffIdle := timeStats[0].Idle - u.allCpuPrevIdle
	diffTotal := total - u.allCpuPrevTotal
	all := common.Rounding(100*(diffTotal-diffIdle)/diffTotal, 2)
	if v := u.loadMin.Value(); v == 0 || all < v {
		u.loadMin.Update(all)
	}
	if v := u.loadMax.Value(); v == 0 || all > v {
		u.loadMax.Update(all)
	}

	u.all.Update(all)
	u.allCpuPrevTotal = total
	u.allCpuPrevIdle = timeStats[0].Idle

	u.AttrMu.Lock()
	u.Attrs[AttrCpuCores].Value = u.cores
	u.Attrs[AttrCpuMhz].Value = u.mhz
	u.Attrs[AttrCpuAll].Value = u.all.Value()
	u.Attrs[AttrLoadMax].Value = u.loadMax.Value()
	u.Attrs[AttrLoadMin].Value = u.loadMin.Value()
	u.AttrMu.Unlock()

	//u.SetMetric(u.Id, "cpuspeed", map[string]float32{
	//	"all": common.Rounding32(u.all.Value(), 2),
	//})

	u.eventBus.Publish(bus.TopicEntities, events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
