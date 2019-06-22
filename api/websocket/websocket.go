package websocket

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/api/websocket/controllers"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/stream"
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
	graceful *graceful_service.GracefulService) *WebSocket {

	server := &WebSocket{
		Controllers: NewControllers(adaptors, stream, endpoint),
	}

	graceful.Subscribe(server)

	return server
}

func (s *WebSocket) Start() {
	log.Infof("Serving websocket service")
}

func (s *WebSocket) Shutdown() {
	log.Info("Server exiting")
}
