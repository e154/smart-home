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
	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/cpu"
	"go.uber.org/atomic"
	"sync"
	"time"
)

type Cpu struct {
	Cores int64   `json:"cores"`
	Mhz   float64 `json:"mhz"`
	All   float64 `json:"all"`
}

type CpuManager struct {
	publisher       IPublisher
	cores           int64
	model           string
	mhz             float64
	all             metrics.GaugeFloat64
	allCpuPrevTotal float64
	allCpuPrevIdle  float64
	quit            chan struct{}
	isStarted       atomic.Bool
	updateLock      sync.Mutex
}

func NewCpuManager(publisher IPublisher) (c *CpuManager) {
	c = &CpuManager{
		all:       metrics.NewGaugeFloat64(),
		quit:      make(chan struct{}),
		publisher: publisher,
	}
	cpuInfo, err := cpu.Info()
	if err != nil {
		return
	}
	c.mhz = cpuInfo[0].Mhz
	c.cores = int64(cpuInfo[0].Cores)
	c.model = cpuInfo[0].Model
	return
}

func (c *CpuManager) start(pause int) {
	if c.isStarted.Load() {
		return
	}
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(pause))
		defer ticker.Stop()

		c.isStarted.Store(true)
		defer func() {
			c.isStarted.Store(false)
		}()

		for {
			select {
			case <-c.quit:
				return
			case <-ticker.C:
				c.selfUpdate()
			}
		}
	}()
}

func (c *CpuManager) stop() {
	if !c.isStarted.Load() {
		return
	}
	c.quit <- struct{}{}
}

func (c *CpuManager) Snapshot() Cpu {
	c.updateLock.Lock()
	defer c.updateLock.Unlock()

	return Cpu{
		Cores: c.cores,
		Mhz:   c.mhz,
		All:   c.all.Value(),
	}
}

func (c *CpuManager) All() float64 {
	return c.all.Value()
}

func (c *CpuManager) selfUpdate() {
	c.updateLock.Lock()
	defer c.updateLock.Unlock()

	timeStats, err := cpu.Times(false)
	if err != nil || len(timeStats) == 0 {
		return
	}

	total := timeStats[0].Total()
	diffIdle := float64(timeStats[0].Idle - c.allCpuPrevIdle)
	diffTotal := float64(total - c.allCpuPrevTotal)
	c.all.Update(100 * (diffTotal - diffIdle) / diffTotal)
	c.allCpuPrevTotal = total
	c.allCpuPrevIdle = timeStats[0].Idle

	c.broadcast()
}

func (c *CpuManager) broadcast() {
	go c.publisher.Broadcast("cpu")
}
