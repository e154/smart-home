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

package hdd

import (
	"sync"

	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/v4/disk"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	cores           int64
	model           string
	total           metrics.Gauge
	free            metrics.Gauge
	usedPercent     metrics.GaugeFloat64
	allCpuPrevTotal float64
	allCpuPrevIdle  float64
	updateLock      *sync.Mutex
	MountPoint      string
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	var mountPoint string
	if _mountPoint, ok := entity.Settings[AttrMountPoint]; ok {
		mountPoint = _mountPoint.String()
	}
	if mountPoint == "" {
		mountPoint = "/"
	}

	actor := &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		updateLock: &sync.Mutex{},
		MountPoint: mountPoint,
	}

	actor.Attrs = NewAttr()
	actor.Setts = entity.Settings

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	if actor.Actions == nil {
		actor.Actions = NewActions()
	}

	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) runAction(_ events.EventCallEntityAction) {
	go e.selfUpdate()
}

func (e *Actor) selfUpdate() {

	e.updateLock.Lock()
	defer e.updateLock.Unlock()

	var mountPoint = "/"
	if e.MountPoint != "" {
		mountPoint = e.MountPoint
	}

	if r, err := disk.Usage(mountPoint); err == nil {
		e.AttrMu.Lock()
		e.Attrs[AttrFstype].Value = r.Fstype
		e.Attrs[AttrTotal].Value = r.Total
		e.Attrs[AttrFree].Value = r.Free
		e.Attrs[AttrUsed].Value = r.Used
		e.Attrs[AttrUsedPercent].Value = r.UsedPercent
		e.Attrs[AttrInodesTotal].Value = r.InodesTotal
		e.Attrs[AttrInodesUsed].Value = r.InodesUsed
		e.Attrs[AttrInodesFree].Value = r.InodesFree
		e.Attrs[AttrInodesUsedPercent].Value = r.InodesUsedPercent
		e.AttrMu.Unlock()
	}
	e.SaveState(false, false)
}
