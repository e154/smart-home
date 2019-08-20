package gate_client

import "net/http"

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

type StreamRequestModel struct {
	URI    string      `json:"uri"`
	Method string      `json:"method"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}

type StreamResponseModel struct {
	Code   int         `json:"code"`
	Body   []byte      `json:"body"`
	Header http.Header `json:"header"`
}
