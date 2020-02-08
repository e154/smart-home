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

package dashboard_models

import (
	"github.com/shirou/gopsutil/mem"
	"sync"
)

type Memory struct {
	sync.Mutex
	SwapTotal uint64 `json:"swap_total"`
	SwapFree  uint64 `json:"swap_free"`
	MemTotal  uint64 `json:"mem_total"`
	MemFree   uint64 `json:"mem_free"`
}

func (m *Memory) Update() {

	m.Lock()
	swap, _ := mem.SwapMemory()
	m.SwapFree = swap.Free / 1024
	m.SwapTotal = swap.Total / 1024

	memory, _ := mem.VirtualMemory()
	m.MemFree = memory.Free / 1024
	m.MemTotal = memory.Total / 1024
	m.Unlock()
}
