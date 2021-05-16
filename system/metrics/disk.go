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

//+build linux windows darwin,!386

package metrics

import (
	"github.com/shirou/gopsutil/disk"
	"go.uber.org/atomic"
	"sync"
	"time"
)

// DiskManager ...
type DiskManager struct {
	publisher  IPublisher
	root       UsageStat
	quit       chan struct{}
	isStarted  atomic.Bool
	updateLock sync.Mutex
}

// NewDiskManager ...
func NewDiskManager(publisher IPublisher) (manager *DiskManager) {
	manager = &DiskManager{
		quit:      make(chan struct{}),
		publisher: publisher,
	}
	manager.selfUpdate()
	return
}

func (d *DiskManager) start(pause int) {
	if d.isStarted.Load() {
		return
	}
	go func() {

		ticker := time.NewTicker(time.Second * time.Duration(pause))
		defer ticker.Stop()

		d.isStarted.Store(true)
		defer func() {
			d.isStarted.Store(false)
		}()

		for {
			select {
			case <-d.quit:

				return
			case <-ticker.C:
				d.selfUpdate()
			}
		}
	}()
}

func (d *DiskManager) stop() {
	if !d.isStarted.Load() {
		return
	}
	d.quit <- struct{}{}
}

// Snapshot ...
func (d *DiskManager) Snapshot() Disk {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	return Disk{
		Root: d.root,
	}
}

func (d *DiskManager) selfUpdate() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	if root, err := disk.Usage("/"); err == nil {
		d.root = UsageStat{
			Path:              root.Path,
			Fstype:            root.Fstype,
			Total:             root.Total,
			Free:              root.Free,
			Used:              root.Used,
			UsedPercent:       root.UsedPercent,
			InodesTotal:       root.InodesTotal,
			InodesUsed:        root.InodesUsed,
			InodesFree:        root.InodesFree,
			InodesUsedPercent: root.InodesUsedPercent,
		}
	}

	d.broadcast()
}

func (d *DiskManager) broadcast() {
	go d.publisher.Broadcast("disk")
}
