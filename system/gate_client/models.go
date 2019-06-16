package gate_client

import "github.com/e154/smart-home/system/uuid"

type Settings struct {
	Id      uuid.UUID `json:"id"`
	Address string    `json:"address"`
	Enabled bool      `json:"enabled"`
}

func (s Settings) Valid() bool {
	if s.Address != "" && s.Enabled {
		return true
	}
	return false
}
