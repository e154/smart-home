// +build windows

package telemetry

//http://markelov.blogspot.ru/2009/01/linux-procmeminfo.html
type Memory struct {
	SwapTotal	uint64		`json:"swap_total"`
	SwapFree	uint64		`json:"swap_free"`
	MemTotal	uint64		`json:"mem_total"`
	MemFree		uint64		`json:"mem_free"`
}

func (m *Memory) Update() {


}
