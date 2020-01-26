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
