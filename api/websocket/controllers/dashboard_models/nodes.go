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
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
)

// Nodes ...
type Nodes struct {
	metric *metrics.MetricManager
}

// NewNode ...
func NewNode(metric *metrics.MetricManager) (node *Nodes) {
	node = &Nodes{metric: metric}
	return
}

// Broadcast ...
func (n *Nodes) Broadcast() (map[string]interface{}, bool) {

	return map[string]interface{}{
		"nodes": n.metric.Node.Snapshot(),
	}, true
}

// only on request: 'dashboard.get.nodes.status'
//
func (n *Nodes) NodesStatus(client stream.IStreamClient, message stream.Message) {

	payload := map[string]interface{}{"nodes": n.metric.Node.Snapshot()}
	response := message.Response(payload)
	client.Write(response.Pack())
}
