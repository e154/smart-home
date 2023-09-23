//go:build (linux && !mips64 && !mips64le) || darwin
// +build linux,!mips64,!mips64le darwin

package uptime

import "github.com/shirou/gopsutil/v3/host"

func GetUptime() (uint64, error) {
	return host.Uptime()
}
