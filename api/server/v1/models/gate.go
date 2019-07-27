package models

// swagger:model
type GateSettings struct {
	GateServerToken string `json:"gate_server_token"`
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
}

// swagger:model
type UpdateGateSettings struct {
	GateServerToken string `json:"gate_server_token"`
	Address         string `json:"address"`
	Enabled         bool   `json:"enabled"`
}

// swagger:model
type GateMobileList struct {
	Total     int64    `json:"total"`
	TokenList []string `json:"token_list"`
}

// swagger:model
type DeleteGateMobile struct {
	Token string `json:"token"`
}
