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

package controllers

import (
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
)

type Workflow struct {
	metric *metrics.MetricManager
}

func NewWorkflow(metric *metrics.MetricManager) *Workflow {
	return &Workflow{
		metric: metric,
	}
}

func (d *Workflow) Broadcast() (map[string]interface{}, bool) {

	log.Infof("broadcast workflow status")

	snapshot := d.metric.Workflow.Snapshot()

	return map[string]interface{}{
		"workflows": snapshot,
	}, true
}

// only on request: 'workflow.get.status'
//
func (d *Workflow) GetWorkflowStatus(client stream.IStreamClient, message stream.Message) {

	log.Infof("get workflow status")

	snapshot := d.metric.Workflow.Snapshot()

	msg := stream.Message{
		Id:      message.Id,
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"workflows": snapshot,
		},
	}

	client.Write(msg.Pack())
}
