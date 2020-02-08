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
	"github.com/shirou/gopsutil/cpu"
	"sync"
)

type Processors struct {
	Processors	[]cpu.InfoStat			`json:"processors"`
}

func NewCpu() (_cpu *Cpu) {

	cpuinfo, err := cpu.Info()
	if err != nil || len(cpuinfo) == 0 {
		return
	}

	_cpu = &Cpu{
		Cpuinfo:			&Processors{
								Processors: cpuinfo,
							},
		Usage:				make(map[int]float64),
		cpu_prev_total:		make(map[int]uint64),
		cpu_prev_idle:		make(map[int]uint64),
	}

	return
}

type Cpu struct {
	sync.Mutex
	Cpuinfo				*Processors 		`json:"processors"`
	All					float64				`json:"all"`
	Usage				map[int]float64		`json:"usage"`
	cpu_prev_total		map[int]uint64
	cpu_prev_idle		map[int]uint64
	all_cpu_prev_total	float64
	all_cpu_prev_idle	float64
}

func (m *Cpu) Update() {

	cpuinfo, err := cpu.Info()
	if err != nil || len(cpuinfo) == 0 {
		return
	}

	m.Lock()
	timeStats, _ := cpu.Times(false)
	total := timeStats[0].Total()
	diff_idle := float64(timeStats[0].Idle - m.all_cpu_prev_idle)
	diff_total := float64(total - m.all_cpu_prev_total)
	m.All = 100 * (diff_total - diff_idle) / diff_total
	m.all_cpu_prev_total = total
	m.all_cpu_prev_idle = timeStats[0].Idle
	m.Unlock()
}
