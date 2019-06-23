package gate_client

type Settings struct {
	GateServerToken string `json:"gate_server_token"`
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
}

func (s Settings) Valid() bool {
	if s.Address != "" && s.Enabled {
		return true
	}
	return false
}

type MobileList struct {
	Total     int64    `json:"total"`
	TokenList []string `json:"token_list"`
}
