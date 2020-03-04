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
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"sync"
	"time"
)

type Nodes struct {
	sync.Mutex
	Total      int64              `json:"total"`
	Status     map[int64]string   `json:"status"`
	lastUpdate time.Time          `json:"-"`
	adaptors   *adaptors.Adaptors `json:"-"`
	core       *core.Core         `json:"-"`
}

func NewNode(adaptors *adaptors.Adaptors,
	core *core.Core) (node *Nodes) {

	node = &Nodes{
		Status:   make(map[int64]string),
		adaptors: adaptors,
		core:     core,
	}

	return
}

func (n *Nodes) Update() {

	if n.core == nil {
		return
	}

	if time.Now().Sub(n.lastUpdate).Seconds() < 15 {
		return
	}

	_, total, _ := n.adaptors.Node.List(999, 0, "", "")

	nodes := n.core.GetNodes()

	n.Lock()

	for key, _ := range n.Status {
		delete(n.Status, key)
	}

	for _, node := range nodes {
		n.Status[node.Id] = node.GetConnStatus()
	}

	n.lastUpdate = time.Now()
	n.Total = total
	n.Unlock()
}

func (n *Nodes) Broadcast() (map[string]interface{}, bool) {

	n.Update()

	return map[string]interface{}{
		"nodes": n,
	}, true
}

// only on request: 'dashboard.get.nodes.status'
//
func (n *Nodes) NodesStatus(client stream.IStreamClient, message stream.Message) {

	n.Update()

	payload := map[string]interface{}{"nodes": n,}
	response := message.Response(payload)
	client.Write(response.Pack())
}
