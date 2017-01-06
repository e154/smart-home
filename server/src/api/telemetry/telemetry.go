package telemetry

import (
	"time"
	"../stream"
	"../log"
	"encoding/json"
)

const (
	SLEEP  = 3
)

var (
	Hub		stream.Hub
	Telemetry        *telemetry = nil
)

type telemetry struct {
	quit chan bool
	memory	*Memory
	cpu	*Cpu
	uptime	*Uptime
}

func (t *telemetry) Run()  {

	for  {
		select {
		case <-t.quit:
			break
		default:

		}

		t.Update()

		time.Sleep(time.Second * SLEEP)
	}
}

func (t *telemetry) Stop() {
	t.quit <- true
}

func (t *telemetry) Update() {

	t.memory.Update()
	t.cpu.Update()
	t.uptime.Update()

	msg, _ := json.Marshal(
		map[string]interface{}{"type": "broadcast",
			"value": map[string]interface{}{"type": "telemetry", "body": map[string]interface{}{
				"memory": t.memory,
				"cpu": t.cpu,
				"time": time.Now(),
				"uptime": t.uptime,
			}}},
	)

	Hub.Broadcast(string(msg))
}

func Initialize() {
	log.Info("Telemetry initialize...")

	if Telemetry != nil {
		return
	}

	Telemetry = &telemetry{
		cpu: NewCpu(),
		memory:	&Memory{},
		uptime: &Uptime{},
		quit: make(chan bool),
	}

	Hub = stream.GetHub()

	go Telemetry.Run()
}