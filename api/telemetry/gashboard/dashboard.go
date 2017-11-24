package dasboard

import (
	"time"
	"reflect"
	"encoding/json"
	"github.com/e154/smart-home/api/stream"
	"github.com/astaxie/beego"
	"github.com/e154/smart-home/api/models"
)

var (
	telemetry_time int = 0
	Hub		stream.Hub
)

func NewDashboard() *Dashboard {

	Hub	= stream.GetHub()

	var err error
	if telemetry_time, err =  beego.AppConfig.Int("telemetry_time"); err != nil {
		telemetry_time = 3
	}

	dashboard := &Dashboard{
		Cpu: NewCpu(),
		Memory:	&Memory{},
		Uptime: &Uptime{},
		Disk: NewDisk(),
		Nodes: NewNode(),
		Devices: &Devices{
			Status: make(map[int64]*models.DeviceState),
		},
		quit: make(chan bool),
	}

	Hub.Subscribe("dashboard.get.nodes.status", dashboard.Nodes.streamNodesStatus)
	Hub.Subscribe("dashboard.get.flows.status", streamFlowsStatus)
	Hub.Subscribe("dashboard.get.telemetry", dashboard.streamTelemetry)

	return dashboard
}

type Dashboard struct {
	quit    chan bool
	Memory  *Memory
	Cpu     *Cpu
	Uptime  *Uptime
	Disk    *Disk
	Nodes   *Nodes
	Devices *Devices
}

func (t *Dashboard) Run()  {

	go func() {
		for  {
			select {
			case <-t.quit:
				break
			default:

			}

			t.broadcastAll()

			time.Sleep(time.Second * time.Duration(telemetry_time))
		}
	}()
}

func (t *Dashboard) Stop() {
	t.quit <- true
}

func (t *Dashboard) BroadcastOne(pack string, id int64) {

	var body interface{}
	var ok bool

	switch pack {
	case "devices":
		body, ok = t.Devices.BroadcastOne(id)
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
		body, ok = t.Devices.Broadcast()
	}

	if (ok) {
		t.sendMsg(body)
	}
}

func (t *Dashboard) sendMsg (body interface{}) {

	msg, _ := json.Marshal(map[string]interface{}{
		"type": "broadcast",
		"value": map[string]interface{}{
			"type": "dashboard.telemetry",
			"body": body,
		},
	})

	Hub.Broadcast(string(msg))
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
		"cpu": map[string]interface{}{"usage":t.Cpu.Usage, "all": t.Cpu.All},
		"time": time.Now(),
		"uptime": t.Uptime,
	})
}

func (t *Dashboard) GetStates() *Dashboard {

	t.Memory.Update()
	t.Cpu.Update()
	t.Uptime.Update()
	t.Disk.Update()
	t.Nodes.Update()
	t.Devices.Update()

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
		"memory": states.Memory,
		"cpu": map[string]interface{}{"usage":t.Cpu.Usage, "info": t.Cpu.Cpuinfo, "all": t.Cpu.All},
		"time": time.Now(),
		"uptime": states.Uptime,
		"disk": states.Disk,
		"nodes": states.Nodes,
		"devices": states.Devices,
	}})
	client.Send(string(msg))
}
