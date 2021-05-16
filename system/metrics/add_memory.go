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

package metrics

import (
	"go.uber.org/atomic"
	"runtime"
	"sync"
	"time"
)

// AppMemory ...
type AppMemory struct {
	Alloc      uint64    `json:"alloc"`
	HeapAlloc  uint64    `json:"heap_alloc"`
	TotalAlloc uint64    `json:"total_alloc"`
	Sys        uint64    `json:"sys"`
	NumGC      uint64    `json:"num_gc"`
	LastGC     time.Time `json:"last_gc"`
}

// AppMemoryManager ...
type AppMemoryManager struct {
	publisher  IPublisher
	isStarted  *atomic.Bool
	quit       chan struct{}
	updateLock *sync.Mutex
	alloc      uint64
	heapAlloc  uint64
	totalAlloc uint64
	sys        uint64
	numGC      uint64
	lastGC     time.Time
}

// NewAppMemoryManager ...
func NewAppMemoryManager(publisher IPublisher) *AppMemoryManager {
	return &AppMemoryManager{
		publisher:  publisher,
		quit:       make(chan struct{}),
		isStarted:  atomic.NewBool(false),
		updateLock: &sync.Mutex{},
	}
}

func (d *AppMemoryManager) start(pause int) {
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

func (d *AppMemoryManager) stop() {
	if !d.isStarted.Load() {
		return
	}
	d.quit <- struct{}{}
}

func (d *AppMemoryManager) selfUpdate() {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	var s runtime.MemStats
	runtime.ReadMemStats(&s)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	d.alloc = bToMb(s.Alloc)
	d.heapAlloc = bToMb(s.HeapAlloc)
	d.totalAlloc = bToMb(s.TotalAlloc)
	d.sys = bToMb(s.Sys)
	d.numGC = uint64(s.NumGC)
	d.lastGC = time.Unix(0, int64(s.LastGC))

	d.broadcast()
}

// Snapshot ...
func (d AppMemoryManager) Snapshot() AppMemory {
	d.updateLock.Lock()
	defer d.updateLock.Unlock()

	return AppMemory{
		Alloc:      d.alloc,
		HeapAlloc:  d.heapAlloc,
		TotalAlloc: d.totalAlloc,
		Sys:        d.sys,
		NumGC:      d.numGC,
		LastGC:     d.lastGC,
	}
}

func (d *AppMemoryManager) broadcast() {
	go d.publisher.Broadcast("app_memory")
}
