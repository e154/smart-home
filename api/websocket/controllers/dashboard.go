package controllers

import (
	dashboardModel "github.com/e154/smart-home/api/websocket/controllers/dashboard_models"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry"
	"time"
)

type ControllerDashboard struct {
	*ControllerCommon
	quit          chan bool
	Nodes         *dashboardModel.Nodes
	Gate          *dashboardModel.Gate
	telemetryTime int
	Memory        *dashboardModel.Memory
	AppMemory     *dashboardModel.AppMemory
	Cpu           *dashboardModel.Cpu
	Uptime        *dashboardModel.Uptime
	Disk          *dashboardModel.Disk
	devices       *dashboardModel.Devices
	Workflow      *dashboardModel.Workflow
}

func NewControllerDashboard(common *ControllerCommon) *ControllerDashboard {
	return &ControllerDashboard{
		ControllerCommon: common,
		telemetryTime:    3,
		quit:             make(chan bool),
		Cpu:              dashboardModel.NewCpu(),
		Memory:           &dashboardModel.Memory{},
		AppMemory:        &dashboardModel.AppMemory{},
		Uptime:           &dashboardModel.Uptime{},
		Disk:             dashboardModel.NewDisk(),
		Nodes:            dashboardModel.NewNode(common.adaptors, common.core),
		devices:          dashboardModel.NewDevices(common.adaptors, common.core),
		Gate:             dashboardModel.NewGate(common.gate),
		Workflow:         dashboardModel.NewWorkflow(common.adaptors, common.core),
	}
}

func (c *ControllerDashboard) Start() {
	c.telemetry.Subscribe("dashboard", c)
	c.stream.Subscribe("dashboard.get.nodes.status", c.Nodes.NodesStatus)
	c.stream.Subscribe("dashboard.get.gate.status", c.Gate.GatesStatus)
	c.stream.Subscribe("t.get.flows.status", dashboardModel.FlowsStatus)
	c.stream.Subscribe("dashboard.get.telemetry", c.Telemetry)

	//i := 0

	go func() {
		for {
			select {
			case <-c.quit:
				break
			default:

			}

			//if i >= 10 {
			//	runtime.GC()
			//	i = 0
			//}

			c.broadcastAll()

			time.Sleep(time.Second * time.Duration(c.telemetryTime))
			//i++
		}
	}()
}

func (c *ControllerDashboard) Stop() {
	c.telemetry.UnSubscribe("dashboard")
	c.stream.UnSubscribe("dashboard.get.nodes.status")
	c.stream.UnSubscribe("dashboard.get.gate.status")
	c.stream.UnSubscribe("t.get.flows.status")
	c.stream.UnSubscribe("dashboard.get.telemetry")

	c.quit <- true
}

func (t *ControllerDashboard) BroadcastOne(param interface{}) {

	var body map[string]interface{}
	var ok bool

	switch v := param.(type) {
	case telemetry.WorkflowScenario:
		body, ok = t.Workflow.BroadcastOne(v)
	case telemetry.Device:
		body, ok = t.devices.BroadcastOne(v.Id, v.ElementName)
	}

	if ok {
		t.sendMsg(body)
	}
}

func (t *ControllerDashboard) Broadcast(param interface{}) {

	var body map[string]interface{}
	var ok bool

	switch param.(type) {
	case telemetry.Node:
		body, ok = t.Nodes.Broadcast()
	case telemetry.Device:
		body, ok = t.devices.Broadcast()
	case telemetry.WorkflowScenario:
		body, ok = t.Workflow.Broadcast()
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

// every time send:
// memory, swap, cpu, uptime
//
func (t *ControllerDashboard) broadcastAll() {

	t.Memory.Update()
	t.AppMemory.Update()
	t.Cpu.Update()
	t.Uptime.Update()
	t.Gate.Update()

	//fmt.Println(t.AppMemory)

	msg := &stream.Message{
		Command: "dashboard.telemetry",
		Type:    stream.Broadcast,
		Forward: stream.Request,
		Payload: map[string]interface{}{
			"memory":     t.Memory,
			"app_memory": t.AppMemory,
			"cpu":        map[string]interface{}{"usage": t.Cpu.Usage, "all": t.Cpu.All},
			"uptime":     t.Uptime,
			"gate":       t.Gate,
		},
	}

	t.stream.Broadcast(msg.Pack())
}

func (t *ControllerDashboard) GetStates() *ControllerDashboard {

	t.Memory.Update()
	t.Cpu.Update()
	t.Uptime.Update()
	t.Disk.Update()
	t.Nodes.Update()
	t.devices.Update()
	t.Gate.Update()

	return t
}

// only on request: 'dashboard.get.telemetry'
//
func (t *ControllerDashboard) Telemetry(client *stream.Client, message stream.Message) {

	states := t.GetStates()
	msg := &stream.Message{
		Id:      message.Id,
		Command: "dashboard.telemetry",
		Forward: stream.Response,
		Payload: map[string]interface{}{
			"memory":  states.Memory,
			"cpu":     map[string]interface{}{"usage": t.Cpu.Usage, "info": t.Cpu.Cpuinfo, "all": t.Cpu.All},
			"time":    time.Now(),
			"uptime":  states.Uptime,
			"disk":    states.Disk,
			"nodes":   states.Nodes,
			"devices": states.devices,
			"gate":    states.Gate,
		},
	}

	client.Send <- msg.Pack()
}
