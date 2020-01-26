// +build linux,!mips64,!mips64le darwin windows

package dashboard_models

import (
	"github.com/shirou/gopsutil/host"
	"sync"
)

type Uptime struct {
	sync.Mutex
	Total uint64 `json:"total"`
	Idle  uint64 `json:"idle"`
}

func (u *Uptime) Update() {

	total, _ := host.Uptime()
	u.Lock()
	u.Total = total
	u.Unlock()
}
