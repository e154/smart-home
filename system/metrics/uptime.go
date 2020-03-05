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

// +build linux,!mips64,!mips64le darwin

package metrics

import (
	"github.com/shirou/gopsutil/host"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type UptimeManager struct {
	updateLock sync.Mutex
	total      uint64
	idle       uint64
	isStarted  atomic.Bool
	quit       chan struct{}
}

func NewUptimeManager() (uptime *UptimeManager) {
	uptime = &UptimeManager{
		quit: make(chan struct{}),
	}
	uptime.selfUpdate()
	return
}

func (d *UptimeManager) Start(pause int) {
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

func (d *UptimeManager) Stop() {
	if !d.isStarted.Load() {
		return
	}
	d.quit <- struct{}{}
}

func (d *UptimeManager) Snapshot() Uptime {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	return Uptime{
		Total: d.total,
	}
}

func (d *UptimeManager) selfUpdate() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	total, _ := host.Uptime()
	d.total = total
}
