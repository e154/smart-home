package dashboard

import (
	"encoding/json"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"github.com/op/go-logging"
	"reflect"
	"time"
)

var (
	log = logging.MustGetLogger("dashboard")
)

type Dashboard struct {
	quit          chan bool
	Nodes         *Nodes
	telemetryTime int
	stream        *stream.StreamService
	Memory        *Memory
	Cpu           *Cpu
	Uptime        *Uptime
	Disk          *Disk
	devices       *Devices
}

func NewDashboard(stream *stream.StreamService,
	adaptors *adaptors.Adaptors) *Dashboard {

	dashboard := &Dashboard{
		Cpu:           NewCpu(),
		Memory:        &Memory{},
		Uptime:        &Uptime{},
		Disk:          NewDisk(),
		Nodes:         NewNode(adaptors),
		quit:          make(chan bool),
		stream:        stream,
		telemetryTime: 3,
		devices:       NewDevices(adaptors),
	}

	return dashboard
}

func (t *Dashboard) Run(core *core.Core) {

	t.devices.core = core
	t.Nodes.core = core

	t.stream.Subscribe("dashboard.get.nodes.status", t.Nodes.streamNodesStatus)
	t.stream.Subscribe("t.get.flows.status", streamFlowsStatus)
	t.stream.Subscribe("dashboard.get.telemetry", t.streamTelemetry)

	go func() {
		for {
			select {
			case <-t.quit:
				break
			default:

			}

			t.broadcastAll()

			time.Sleep(time.Second * time.Duration(t.telemetryTime))
		}
	}()
}

func (t *Dashboard) Stop() {

	t.stream.UnSubscribe("dashboard.get.nodes.status")
	t.stream.UnSubscribe("t.get.flows.status")
	t.stream.UnSubscribe("dashboard.get.telemetry")

	t.quit <- true
}

func (t *Dashboard) BroadcastOne(pack string, deviceId int64, elementName string) {

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

func (t *Dashboard) Broadcast(pack string) {

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

func (t *Dashboard) sendMsg(body interface{}) {

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
func (t *Dashboard) broadcastAll() {

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

func (t *Dashboard) GetStates() *Dashboard {

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
func (t *Dashboard) streamTelemetry(client *stream.Client, value interface{}) {
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
