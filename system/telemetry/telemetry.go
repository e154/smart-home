package telemetry

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry/dashboard"
	"github.com/e154/smart-home/system/telemetry/telemetry_map"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("telemetry")
)


type Telemetry struct {
	dashboard *dashboard.Dashboard
	_map      *telemetry_map.Map
	core      *core.Core
}

func NewTelemetry(
	stream *stream.StreamService,
	adaptors *adaptors.Adaptors) core.ITelemetry {

	telemetry := &Telemetry{
		dashboard: dashboard.NewDashboard(stream, adaptors),
		_map:      telemetry_map.NewMap(stream, adaptors),
	}

	return telemetry
}

func (t *Telemetry) Run(core *core.Core) {

	log.Info("Run")

	t.dashboard.Run(core)
	t._map.Run(core)
}

func (t *Telemetry) Stop() {

	log.Info("Stop")

	t.dashboard.Stop()
	t._map.Stop()
}

func (t *Telemetry) Broadcast(pack string) {

	t.dashboard.Broadcast(pack)
	t._map.Broadcast(pack)
}

func (t *Telemetry) BroadcastOne(pack string, deviceId int64, elementName string) {

	t.dashboard.BroadcastOne(pack, deviceId, elementName)
	t._map.BroadcastOne(pack, deviceId, elementName)
}
