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
	"fmt"
	"runtime"
	"sync"
	"time"
)

type AppMemory struct {
	sync.Mutex
	Alloc      uint64    `json:"alloc"`
	HeapAlloc  uint64    `json:"heap_alloc"`
	TotalAlloc uint64    `json:"total_alloc"`
	Sys        uint64    `json:"sys"`
	NumGC      uint64    `json:"num_gc"`
	LastGC     time.Time `json:"last_gc"`
}

func (m *AppMemory) Update() {

	var s runtime.MemStats
	runtime.ReadMemStats(&s)

	m.Lock()
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	m.Alloc = bToMb(s.Alloc)
	m.HeapAlloc = bToMb(s.HeapAlloc)
	m.TotalAlloc = bToMb(s.TotalAlloc)
	m.Sys = bToMb(s.Sys)
	m.NumGC = uint64(s.NumGC)
	m.LastGC = time.Unix(0, int64(s.LastGC))
	m.Unlock()
}

func (m *AppMemory) String() string {
	return fmt.Sprintf("Alloc: %v, HeapAlloc: %v, TotalAlloc: %v, Sys: %v, NumGC: %v, LastGC: %v", m.Alloc, m.HeapAlloc, m.TotalAlloc, m.Sys, m.NumGC, m.LastGC)
}
