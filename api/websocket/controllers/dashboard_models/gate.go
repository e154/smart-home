package dashboard_models

import (
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/stream"
)

type Gate struct {
	gate   *gate_client.GateClient
	status string
}

func NewGate(gate *gate_client.GateClient) *Gate {
	return &Gate{
		gate:   gate,
		status: gate_client.GateStatusWait,
	}
}

func (g *Gate) Update() {
	g.status = g.gate.Status()
}

// only on request: 'dashboard.get.gate.status'
//
func (g *Gate) GatesStatus(client *stream.Client, message stream.Message) {

	g.Update()

	payload := map[string]interface{}{"gate_status": g.status}
	response := message.Response(payload)
	client.Send <- response.Pack()

	return
}
