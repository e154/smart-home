package telemetry

import (
	"time"
	"github.com/astaxie/beego"
	"encoding/json"
	"../stream"
	"../log"
	"reflect"
)

var (
	telemetry_time	int
	Hub		stream.Hub
	Telemetry        *telemetry = nil
)

type telemetry struct {
	quit    chan bool
	Memory  *Memory
	Cpu     *Cpu
	Uptime  *Uptime
	Disk    *Disk
	Nodes   *Nodes
	Devices *Devices
}

func (t *telemetry) Run()  {

	for  {
		select {
		case <-t.quit:
			break
		default:

		}

		t.broadcastAll()

		time.Sleep(time.Second * time.Duration(telemetry_time))
	}
}

func (t *telemetry) Stop() {
	t.quit <- true
}

func (t *telemetry) BroadcastOne(pack string, id int64) {
	switch pack {

	case "devices":
		go t.Devices.BroadcastOne(id)
	}
}

func (t *telemetry) Broadcast(pack string) {
	switch pack {
	case "nodes":
		go t.Nodes.Broadcast()
	case "devices":
		go t.Devices.Broadcast()
	}
}

// every time send:
// memory, swap, cpu, uptime
//
func (t *telemetry) broadcastAll() {

	t.Memory.Update()
	t.Cpu.Update()
	t.Uptime.Update()

	msg, _ := json.Marshal(
		map[string]interface{}{"type": "broadcast",
			"value": map[string]interface{}{"type": "telemetry", "body": map[string]interface{}{
				"memory": t.Memory,
				"cpu": map[string]interface{}{"usage":t.Cpu.Usage, "all": t.Cpu.All},
				"time": time.Now(),
				"uptime": t.Uptime,
			}}},
	)

	Hub.Broadcast(string(msg))
}

func (t *telemetry) GetStates() *telemetry {

	t.Memory.Update()
	t.Cpu.Update()
	t.Uptime.Update()
	t.Disk.Update()
	t.Nodes.Update()
	t.Devices.Update()

	return t
}

// only on request: 'get.telemetry'
//
func streamTelemetry(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	t := Telemetry.GetStates()
	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "telemetry":
	map[string]interface{}{
		"memory": t.Memory,
		"cpu": map[string]interface{}{"usage":t.Cpu.Usage, "info": t.Cpu.Cpuinfo, "all": t.Cpu.All},
		"time": time.Now(),
		"uptime": t.Uptime,
		"disk": t.Disk,
		"nodes": t.Nodes,
		"devices": t.Devices,
	}})
	client.Send(string(msg))
}

func Initialize() *telemetry {
	log.Info("Telemetry initialize...")

	if Telemetry != nil {
		return Telemetry
	}

	var err error
	if telemetry_time, err =  beego.AppConfig.Int("telemetry_time"); err != nil {
		telemetry_time = 3
	}

	Telemetry = &telemetry{
		Cpu: NewCpu(),
		Memory:	&Memory{},
		Uptime: &Uptime{},
		Disk: NewDisk(),
		Nodes: NewNode(),
		Devices: &Devices{},
		quit: make(chan bool),
	}

	Hub = stream.GetHub()
	Hub.Subscribe("get.nodes.status", streamNodesStatus)
	Hub.Subscribe("get.flows.status", streamFlowsStatus)
	Hub.Subscribe("get.devices.states", streamGetDevicesStates)
	Hub.Subscribe("get.telemetry", streamTelemetry)

	go Telemetry.Run()

	return Telemetry
}