package websocket

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/api/websocket/controllers"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("websocket")
)

type WebSocket struct {
	Controllers *Controllers
}

func NewWebSocket(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint *endpoint.Endpoint,
	scripts *scripts.ScriptService,
	core *core.Core,
	graceful *graceful_service.GracefulService,
	telemetry *telemetry.Telemetry,
	gate *gate_client.GateClient) *WebSocket {

	server := &WebSocket{
		Controllers: NewControllers(adaptors, stream, scripts, core, endpoint, telemetry, gate),
	}

	graceful.Subscribe(server)

	return server
}

func (s *WebSocket) Start() {
	log.Infof("Serving websocket service")
	s.Controllers.Start()
}

func (s *WebSocket) Shutdown() {
	log.Info("Server exiting")
	s.Controllers.Stop()
}
