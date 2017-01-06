package telemetry

import (
	"time"
	"../stream"
	"../log"
	"encoding/json"
	"github.com/astaxie/beego"
)

var (
	telemetry_time	int
	Hub		stream.Hub
	Telemetry        *telemetry = nil
)

type telemetry struct {
	quit chan bool
	memory	*Memory
	cpu	*Cpu
	uptime	*Uptime
	disk	*Disk
}

func (t *telemetry) Run()  {

	for  {
		select {
		case <-t.quit:
			break
		default:

		}

		t.Update()

		time.Sleep(time.Second * time.Duration(telemetry_time))
	}
}

func (t *telemetry) Stop() {
	t.quit <- true
}

func (t *telemetry) Update() {

	t.memory.Update()
	t.cpu.Update()
	t.uptime.Update()
	t.disk.Update()

	msg, _ := json.Marshal(
		map[string]interface{}{"type": "broadcast",
			"value": map[string]interface{}{"type": "telemetry", "body": map[string]interface{}{
				"memory": t.memory,
				"cpu": t.cpu,
				"time": time.Now(),
				"uptime": t.uptime,
				"disk": t.disk,
			}}},
	)

	Hub.Broadcast(string(msg))
}

func Initialize() {
	log.Info("Telemetry initialize...")

	if Telemetry != nil {
		return
	}

	var err error
	if telemetry_time, err =  beego.AppConfig.Int("telemetry_time"); err != nil {
		telemetry_time = 3
	}

	Telemetry = &telemetry{
		cpu: NewCpu(),
		memory:	&Memory{},
		uptime: &Uptime{},
		disk: NewDisk(),
		quit: make(chan bool),
	}

	Hub = stream.GetHub()

	go Telemetry.Run()
}