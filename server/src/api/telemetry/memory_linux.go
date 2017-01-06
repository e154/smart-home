// +build linux

package telemetry

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"../../lib/common"
)

//http://markelov.blogspot.ru/2009/01/linux-procmeminfo.html
type Memory struct {
	SwapTotal	string		`json:"swap_total"`
	SwapFree	string		`json:"swap_free"`
	MemTotal	string		`json:"mem_total"`
	MemFree		string		`json:"mem_free"`
}

func (m *Memory) Update() {

	mem, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return
	}

	m.SwapTotal = common.IBytes(mem.SwapTotal * 1024)
	m.SwapFree = common.IBytes(mem.SwapFree * 1024)
	m.MemTotal = common.IBytes(mem.MemTotal * 1024)
	m.MemFree = common.IBytes(mem.MemFree * 1024)
}
