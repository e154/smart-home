// +build mips64,mips64le

package dasboard

type Uptime struct {
	Total uint64 `json:"total"`
	Idle  uint64 `json:"idle"`
}

func (u *Uptime) Update() (*Uptime, error) {

	return u, nil
}