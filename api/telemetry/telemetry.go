package telemetry

import (
	"github.com/e154/smart-home/api/log"
	"github.com/e154/smart-home/api/telemetry/gashboard"
	"github.com/e154/smart-home/api/telemetry/telemetry_map"
)

var (
	Telemetry		*telemetry = nil
)

type telemetry struct {
	Dashboard	*dasboard.Dashboard
	Map			*telemetry_map.Map
}

func (t *telemetry) Run()  {

	t.Dashboard.Run()
	t.Map.Run()
}

func (t *telemetry) Stop() {

	t.Dashboard.Stop()
	t.Map.Stop()
}

func (t *telemetry) Broadcast(pack string) {

	t.Dashboard.Broadcast(pack)
	t.Map.Broadcast(pack)
}

func (t *telemetry) BroadcastOne(pack string, id int64) {

	t.Dashboard.BroadcastOne(pack, id)
	t.Map.BroadcastOne(pack, id)
}

func Initialize() *telemetry {
	log.Info("Telemetry initialize...")

	if Telemetry != nil {
		return Telemetry
	}

	Telemetry = &telemetry{
		Dashboard: dasboard.NewDashboard(),
		Map: telemetry_map.NewMap(),
	}

	Telemetry.Run()

	return Telemetry
}