// +build linux

package telemetry

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/e154/smart-home/api/log"
)

func NewCpu() (cpu *Cpu) {

	cpuinfo, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		log.Error("telemetry cpu:", err.Error())
	}

	cpu = &Cpu{
		Cpuinfo: 		cpuinfo,
		Usage:			make(map[int]float64),
		cpu_prev_total:		make(map[int]uint64),
		cpu_prev_idle:		make(map[int]uint64),
	}

	return
}

type Cpu struct {
	Cpuinfo			*linuxproc.CPUInfo	`json:"cpuinfo"`
	All			float64			`json:"all"`
	Usage			map[int]float64		`json:"usage"`
	cpu_prev_total		map[int]uint64
	cpu_prev_idle		map[int]uint64
	all_cpu_prev_total	uint64
	all_cpu_prev_idle	uint64
}

func (m *Cpu) Update() {

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return
	}

	s := stat.CPUStatAll
	total := s.User + s.System + s.Nice + s.Idle + s.IOWait + s.IRQ + s.SoftIRQ
	diff_idle := float64(s.Idle - m.all_cpu_prev_idle)
	diff_total := float64(total - m.all_cpu_prev_total)
	m.All = 100 * (diff_total - diff_idle) / diff_total
	m.all_cpu_prev_total = total
	m.all_cpu_prev_idle = s.Idle

	for id, s := range stat.CPUStats {
		total := s.User + s.System + s.Nice + s.Idle + s.IOWait + s.IRQ + s.SoftIRQ
		diff_idle := float64(s.Idle - m.cpu_prev_idle[id])
		diff_total := float64(total - m.cpu_prev_total[id])
		m.Usage[id] = 100 * (diff_total - diff_idle) / diff_total
		m.cpu_prev_total[id] = total
		m.cpu_prev_idle[id] = s.Idle
	}
}