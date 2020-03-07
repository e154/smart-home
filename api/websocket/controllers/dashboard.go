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
	dashboardModel "github.com/e154/smart-home/api/websocket/controllers/dashboard_models"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/stream"
	"time"
)

type ControllerDashboard struct {
	*ControllerCommon
	Nodes    *dashboardModel.Nodes
	devices  *dashboardModel.Devices
	Workflow *dashboardModel.Workflow
	Gate     *dashboardModel.Gate
}

func NewControllerDashboard(common *ControllerCommon) *ControllerDashboard {
	return &ControllerDashboard{
		ControllerCommon: common,
		Nodes:            dashboardModel.NewNode(common.metric),
		devices:          dashboardModel.NewDevices(common.metric),
		Workflow:         dashboardModel.NewWorkflow(common.metric),
		Gate:             dashboardModel.NewGate(common.metric),
	}
}

func (c *ControllerDashboard) Start() {
	c.metric.Subscribe("dashboard", c)
	c.stream.Subscribe("dashboard.get.nodes.status", c.Nodes.NodesStatus)
	c.stream.Subscribe("dashboard.get.gate.status", c.Gate.Status)
	//c.stream.Subscribe("t.get.flows.status", dashboardModel.FlowsStatus)
	c.stream.Subscribe("dashboard.get.telemetry", c.Telemetry)
}

func (c *ControllerDashboard) Stop() {
	c.metric.UnSubscribe("dashboard")
	c.stream.UnSubscribe("dashboard.get.nodes.status")
	c.stream.UnSubscribe("dashboard.get.gate.status")
	//c.stream.UnSubscribe("t.get.flows.status")
	c.stream.UnSubscribe("dashboard.get.telemetry")
}

func (t *ControllerDashboard) Broadcast(param interface{}) {

	var body map[string]interface{}
	var ok bool

	switch v := param.(type) {
	case string:
		switch v {
		case "workflow":
			body, ok = t.Workflow.Broadcast()
		case "node":
			body, ok = t.Nodes.Broadcast()
		case "device":
			body, ok = t.devices.Broadcast()
		}
	case metrics.MapElementCursor:

	}

	if ok {
		t.sendMsg(body)
	}
}

func (t *ControllerDashboard) sendMsg(payload map[string]interface{}) {

	msg := &stream.Message{
		Command: "dashboard.telemetry",
		Type:    stream.Broadcast,
		Forward: stream.Request,
		Payload: payload,
	}

	t.stream.Broadcast(msg.Pack())
}

// only on request: 'dashboard.get.telemetry'
//
func (t *ControllerDashboard) Telemetry(client stream.IStreamClient, message stream.Message) {

	msg := stream.Message{
		Id:      message.Id,
		Command: "dashboard.telemetry",
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"memory":  t.metric.Memory.Snapshot(),
			"cpu":     t.metric.Cpu.Snapshot(),
			"time":    time.Now(),
			"uptime":  t.metric.Uptime.Snapshot(),
			"disk":    t.metric.Disk.Snapshot(),
			"nodes":   t.metric.Node.Snapshot(),
			"devices": t.metric.Device.Snapshot(),
			"gate":    t.metric.Gate.Snapshot(),
		},
	}

	client.Write(msg.Pack())
}
