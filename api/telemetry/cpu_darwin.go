// +build darwin

package telemetry

func NewCpu() (cpu *Cpu) {

	cpu = &Cpu{
		Usage:			make(map[int]float64),
		cpu_prev_total:		make(map[int]uint64),
		cpu_prev_idle:		make(map[int]uint64),
	}

	return
}

type Cpu struct {
	Cpuinfo			string			`json:"cpuinfo"`
	All			float64			`json:"all"`
	Usage			map[int]float64		`json:"usage"`
	cpu_prev_total		map[int]uint64
	cpu_prev_idle		map[int]uint64
	all_cpu_prev_total	uint64
	all_cpu_prev_idle	uint64
}

func (m *Cpu) Update() {

}