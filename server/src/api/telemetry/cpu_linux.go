// +build linux

package telemetry

import (
	linuxproc "github.com/c9s/goprocinfo/linux"
	"../log"
)

func NewCpu() (cpu *Cpu) {

	cpuinfo, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		log.Error("telemetry cpu:", err.Error())
	}

	cpu = &Cpu{
		Cpuinfo: cpuinfo,
	}

	return
}

type Cpu struct {
	Cpuinfo		*linuxproc.CPUInfo	`json:"cpuinfo"`
	Usage		float64			`json:"usage"`
	cpu_prev_total	uint64
	cpu_prev_idle	uint64
}

func (m *Cpu) Update() {

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return
	}

	s := stat.CPUStatAll
	
	total := s.User + s.System + s.Nice + s.Idle + s.IOWait + s.IRQ + s.SoftIRQ
	diff_idle := float64(s.Idle - m.cpu_prev_idle)
	diff_total := float64(total - m.cpu_prev_total)
	m.Usage = 100 * (diff_total - diff_idle) / diff_total
	m.cpu_prev_total = total
	m.cpu_prev_idle = s.Idle
}