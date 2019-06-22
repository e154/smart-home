package controllers

import (
	"encoding/json"
	dashboardModel "github.com/e154/smart-home/api/websocket/controllers/dashboard_models"
	"github.com/e154/smart-home/system/stream"
	"reflect"
	"time"
)

type ControllerDashboard struct {
	*ControllerCommon
	quit          chan bool
	Nodes         *dashboardModel.Nodes
	telemetryTime int
	Memory        *dashboardModel.Memory
	Cpu           *dashboardModel.Cpu
	Uptime        *dashboardModel.Uptime
	Disk          *dashboardModel.Disk
	devices       *dashboardModel.Devices
}

func NewControllerDashboard(common *ControllerCommon) *ControllerDashboard {
	return &ControllerDashboard{
		ControllerCommon: common,
		telemetryTime:    3,
		quit:             make(chan bool),
		Cpu:              dashboardModel.NewCpu(),
		Memory:           &dashboardModel.Memory{},
		Uptime:           &dashboardModel.Uptime{},
		Disk:             dashboardModel.NewDisk(),
		Nodes:            dashboardModel.NewNode(common.adaptors, common.core),
		devices:          dashboardModel.NewDevices(common.adaptors, common.core),
	}
}

func (c *ControllerDashboard) Start() {
	c.telemetry.Subscribe("dashboard", c)
	c.stream.Subscribe("dashboard.get.nodes.status", c.Nodes.NodesStatus)
	c.stream.Subscribe("t.get.flows.status", dashboardModel.FlowsStatus)
	c.stream.Subscribe("dashboard.get.telemetry", c.Telemetry)

	go func() {
		for {
			select {
			case <-c.quit:
				break
			default:

			}

			c.broadcastAll()

			time.Sleep(time.Second * time.Duration(c.telemetryTime))
		}
	}()
}

func (c *ControllerDashboard) Stop() {
	c.telemetry.UnSubscribe("dashboard")
	c.stream.UnSubscribe("dashboard.get.nodes.status")
	c.stream.UnSubscribe("t.get.flows.status")
	c.stream.UnSubscribe("dashboard.get.telemetry")

	c.quit <- true
}

func (t *ControllerDashboard) BroadcastOne(pack string, deviceId int64, elementName string) {

	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = t.devices.BroadcastOne(deviceId, elementName)
	}

	if ok {
		t.sendMsg(body)
	}
}

func (t *ControllerDashboard) Broadcast(pack string) {

	var body interface{}
	var ok bool

	switch pack {
	case "nodes":
		body, ok = t.Nodes.Broadcast()
	case "devices":
		body, ok = t.devices.Broadcast()
	}

	if (ok) {
		t.sendMsg(body)
	}
}

func (t *ControllerDashboard) sendMsg(body interface{}) {

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "broadcast",
		"value": map[string]interface{}{
			"type": "dashboard.telemetry",
			"body": body,
		},
	})

	t.stream.Broadcast(msg)
}

// every time send:
// memory, swap, cpu, uptime
//
func (t *ControllerDashboard) broadcastAll() {

	t.Memory.Update()
	t.Cpu.Update()
	t.Uptime.Update()

	t.sendMsg(map[string]interface{}{
		"memory": t.Memory,
		"cpu":    map[string]interface{}{"usage": t.Cpu.Usage, "all": t.Cpu.All},
		"time":   time.Now(),
		"uptime": t.Uptime,
	})
}

func (t *ControllerDashboard) GetStates() *ControllerDashboard {

	t.Memory.Update()
	t.Cpu.Update()
	t.Uptime.Update()
	t.Disk.Update()
	t.Nodes.Update()
	t.devices.Update()

	return t
}

// only on request: 'dashboard.get.telemetry'
//
func (t *ControllerDashboard) Telemetry(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	states := t.GetStates()
	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "dashboard.telemetry":
	map[string]interface{}{
		"memory":  states.Memory,
		"cpu":     map[string]interface{}{"usage": t.Cpu.Usage, "info": t.Cpu.Cpuinfo, "all": t.Cpu.All},
		"time":    time.Now(),
		"uptime":  states.Uptime,
		"disk":    states.Disk,
		"nodes":   states.Nodes,
		"devices": states.devices,
	}})
	client.Send <- msg
}