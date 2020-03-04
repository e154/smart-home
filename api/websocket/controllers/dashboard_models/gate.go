// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

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
func (g *Gate) GatesStatus(client stream.IStreamClient, message stream.Message) {

	g.Update()

	g.Lock()
	payload := map[string]interface{}{
		"gate_status":  g.Status,
		"access_token": g.AccessToken,
	}
	g.Unlock()

	response := message.Response(payload)
	client.Write(response.Pack())

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
