package dashboard_models

import (
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/stream"
	"sync"
)

type Gate struct {
	gate        *gate_client.GateClient
	sync.Mutex
	Status      string `json:"status"`
	AccessToken string `json:"-"`
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

	status := g.gate.Status()

	g.Lock()
	g.Status = status
	g.AccessToken = settings.GateServerToken
	g.Unlock()
}

// only on request: 'dashboard.get.gate.status'
//
func (g *Gate) GatesStatus(client *stream.Client, message stream.Message) {

	g.Update()

	g.Lock()
	payload := map[string]interface{}{
		"gate_status":  g.Status,
		"access_token": g.AccessToken,
	}
	g.Unlock()

	response := message.Response(payload)
	client.Send <- response.Pack()

	return
}

func (g *Gate) Broadcast() (map[string]interface{}, bool) {

	g.Update()

	g.Lock()
	defer g.Unlock()

	return map[string]interface{}{
		"gate_status": g.Status,
	}, true
}
