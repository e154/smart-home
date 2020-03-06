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

package metrics

import (
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type Memory struct {
	SwapTotal uint64 `json:"swap_total"`
	SwapFree  uint64 `json:"swap_free"`
	MemTotal  uint64 `json:"mem_total"`
	MemFree   uint64 `json:"mem_free"`
}

type MemoryManager struct {
	updateLock sync.Mutex
	swapTotal  uint64
	swapFree   uint64
	memTotal   uint64
	memFree    uint64
	isStarted  atomic.Bool
	quit       chan struct{}
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{}
}

func (d *MemoryManager) start(pause int) {
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

func (d *MemoryManager) stop() {
	if !d.isStarted.Load() {
		return
	}
	d.quit <- struct{}{}
}

func (d *MemoryManager) Snapshot() Memory {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	return Memory{
		SwapTotal: d.swapTotal,
		SwapFree:  d.swapFree,
		MemTotal:  d.memTotal,
		MemFree:   d.memFree,
	}
}

func (d *MemoryManager) selfUpdate() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	swap, _ := mem.SwapMemory()
	d.swapFree = swap.Free / 1024
	d.swapTotal = swap.Total / 1024

	memory, _ := mem.VirtualMemory()
	d.memFree = memory.Free / 1024
	d.memTotal = memory.Total / 1024
}
