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

package hdd

import (
	"fmt"
	"github.com/e154/smart-home/system/event_bus/events"
	"github.com/shirou/gopsutil/v3/disk"
	"sync"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/rcrowley/go-metrics"
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
	eventBus        event_bus.EventBus
	updateLock      *sync.Mutex
	MountPoint      string
}

// NewActor ...
func NewActor(settings m.Attributes,
	entityManager entity_manager.EntityManager,
	eventBus event_bus.EventBus) *Actor {

	var mountPoint string
	if _mountPoint, ok := settings[AttrMountPoint]; ok {
		mountPoint = _mountPoint.String()
	}
	if mountPoint == "" {
		mountPoint = "/"
	}

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
		MountPoint:  mountPoint,
	}

	actor.Manager = entityManager
	actor.Attrs = NewAttr()
	actor.Setts = settings

	if actor.Setts == nil {
		actor.Setts = NewSettings()
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

	var mountPoint = "/"
	if u.MountPoint != "" {
		mountPoint = u.MountPoint
	}

	if r, err := disk.Usage(mountPoint); err == nil {
		u.AttrMu.Lock()
		u.Attrs[AttrFstype].Value = r.Fstype
		u.Attrs[AttrTotal].Value = r.Total
		u.Attrs[AttrFree].Value = r.Free
		u.Attrs[AttrUsed].Value = r.Used
		u.Attrs[AttrUsedPercent].Value = r.UsedPercent
		u.Attrs[AttrInodesTotal].Value = r.InodesTotal
		u.Attrs[AttrInodesUsed].Value = r.InodesUsed
		u.Attrs[AttrInodesFree].Value = r.InodesFree
		u.Attrs[AttrInodesUsedPercent].Value = r.InodesUsedPercent
		u.AttrMu.Unlock()
	}

	u.eventBus.Publish(event_bus.TopicEntities, events.EventStateChanged{
		StorageSave: false,
		PluginName:  u.Id.PluginName(),
		EntityId:    u.Id,
		OldState:    oldState,
		NewState:    u.GetEventState(u),
	})
}
