package dashboard_models

import (
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/stream"
)

type Gate struct {
	gate        *gate_client.GateClient
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
}

func NewGate(gate *gate_client.GateClient) *Gate {
	return &Gate{
		gate:   gate,
		Status: gate_client.GateStatusWait,
	}
}

func (g *Gate) Update() {
	settings, err := g.gate.GetSettings()
	if err != nil {
		return
	}
	g.Status = g.gate.Status()
	g.AccessToken = settings.GateServerToken
}

// only on request: 'dashboard.get.gate.status'
//
func (g *Gate) GatesStatus(client *stream.Client, message stream.Message) {

	payload, _ := g.Broadcast()
	response := message.Response(payload)
	client.Send <- response.Pack()

	return
}

func (g *Gate) Broadcast() (map[string]interface{}, bool) {

	g.Update()

	return map[string]interface{}{
		"gate_status":  g.Status,
		"access_token": g.AccessToken,
	}, true
}
