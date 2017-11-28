// +build !mips64,!mips64le

package dasboard

import (
	"github.com/shirou/gopsutil/host"
)

type Uptime struct {
	Total uint64 `json:"total"`
	Idle  uint64 `json:"idle"`
}

func (u *Uptime) Update() (*Uptime, error) {

	u.Total, _ = host.Uptime()

	return u, nil
}