// +build linux

package dasboard

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
)

//http://markelov.blogspot.ru/2009/01/linux-procmeminfo.html
type Memory struct {
	SwapTotal	uint64		`json:"swap_total"`
	SwapFree	uint64		`json:"swap_free"`
	MemTotal	uint64		`json:"mem_total"`
	MemFree		uint64		`json:"mem_free"`
}

func (m *Memory) Update() {

	mem, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return
	}

	m.SwapTotal = mem.SwapTotal
	m.SwapFree = mem.SwapFree
	m.MemTotal = mem.MemTotal
	m.MemFree = mem.MemFree
}
