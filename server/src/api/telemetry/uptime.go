package telemetry

import (
	"io/ioutil"
	"strings"
	"strconv"
)

type Uptime struct {
	Total float64 `json:"total"`
	Idle  float64 `json:"idle"`
}

func (u *Uptime) Update() (*Uptime, error) {
	b, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(string(b))
	if u.Total, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return nil, err
	}
	if u.Idle, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return nil, err
	}

	return u, nil
}