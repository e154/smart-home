package controllers

import (
	"github.com/e154/smart-home/system/stream"
)

type ControllerDashboard struct {
	*ControllerCommon
}

func NewControllerDashboard(common *ControllerCommon) *ControllerDashboard {
	board := &ControllerDashboard{
		ControllerCommon: common,
	}

	// register methods
	//stream.Subscribe("dashboard.get.nodes.status", board.NodesStatus)
	//stream.Subscribe("t.get.flows.status", board.FlowsStatus)
	//stream.Subscribe("dashboard.get.telemetry", board.Telemetry)


	return board
}

func (n *ControllerDashboard) NodesStatus(client *stream.Client, value interface{}) {
	//v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	//if !ok {
	//	return
	//}
	//
	//n.Update()
	//
	//msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "nodes": n})
	//client.Send <- msg
}

func (n *ControllerDashboard) FlowsStatus(client *stream.Client, value interface{}) {

	return
}

func (t *ControllerDashboard) Telemetry(client *stream.Client, value interface{}) {
	//v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	//if !ok {
	//	return
	//}
	//
	//states := t.GetStates()
	//msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "dashboard.telemetry":
	//map[string]interface{}{
	//	"memory":  states.Memory,
	//	"cpu":     map[string]interface{}{"usage": t.Cpu.Usage, "info": t.Cpu.Cpuinfo, "all": t.Cpu.All},
	//	"time":    time.Now(),
	//	"uptime":  states.Uptime,
	//	"disk":    states.Disk,
	//	"nodes":   states.Nodes,
	//	"devices": states.devices,
	//}})
	//client.Send <- msg
}
