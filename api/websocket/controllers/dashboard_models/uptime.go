// +build linux,!mips64,!mips64le darwin windows

package dashboard_models

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