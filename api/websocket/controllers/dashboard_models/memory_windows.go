// +build windows

package dashboard_models

import (
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
	m.SwapFree = swap.Free / 1024
	m.SwapTotal = swap.Total / 1024
	m.MemFree = memory.Free / 1024
	m.MemTotal = memory.Total / 1024
	m.Unlock()
}
