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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/cpu"
	"sync"
)

type EntityActor struct {
	entity_manager.BaseActor
	cores           int64
	model           string
	mhz             float64
	all             metrics.GaugeFloat64
	allCpuPrevTotal float64
	allCpuPrevIdle  float64
	updateLock      *sync.Mutex
}

func NewEntityActor(entityManager entity_manager.EntityManager) *EntityActor {

	actor := &EntityActor{
		BaseActor: entity_manager.BaseActor{
			Id:                common.EntityId(fmt.Sprintf("%s.%s", EntityCpuspeed, Name)),
			Name:              Name,
			EntityType:        EntityCpuspeed,
			UnitOfMeasurement: "GHz",
			AttrMu:            &sync.Mutex{},
			Attrs:             NewAttr(),
			Manager:           entityManager,
		},
		all:        metrics.NewGaugeFloat64(),
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

func (e *EntityActor) Spawn() entity_manager.PluginActor {
	return e
}

func (u *EntityActor) selfUpdate() {

	u.updateLock.Lock()
	defer u.updateLock.Unlock()

	oldState := u.GetEventState(u)
	u.Now(oldState)

	timeStats, err := cpu.Times(false)
	if err != nil || len(timeStats) == 0 {
		return
	}

	total := timeStats[0].Total()
	diffIdle := timeStats[0].Idle - u.allCpuPrevIdle
	diffTotal := total - u.allCpuPrevTotal
	u.all.Update(100 * (diffTotal - diffIdle) / diffTotal)
	u.allCpuPrevTotal = total
	u.allCpuPrevIdle = timeStats[0].Idle

	u.AttrMu.Lock()
	u.Attrs[AttrCpuCores].Value = u.cores
	u.Attrs[AttrCpuMhz].Value = u.mhz
	u.Attrs[AttrCpuAll].Value = u.all.Value()
	u.AttrMu.Unlock()

	u.Manager.SetMetric(u.Id, "cpuspeed", map[string]interface{}{
		"all": common.Rounding(u.all.Value(), 2),
	})

	u.Send(entity_manager.MessageStateChanged{
		StorageSave: false,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
