package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("stream.controllers")
)

type Controllers struct {
	Image     *ControllerImage
	Worker    *ControllerWorker
	Action    *ControllerAction
	Dashboard *ControllerDashboard
	Map       *ControllerMap
}

func NewControllers(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	scripts *scripts.ScriptService,
	core *core.Core,
	endpoint *endpoint.Endpoint,
	telemetry *telemetry.Telemetry,
	gate *gate_client.GateClient) *Controllers {
	common := NewControllerCommon(adaptors, stream, endpoint, scripts, core, telemetry, gate)
	return &Controllers{
		Image:     NewControllerImage(common),
		Worker:    NewControllerWorker(common),
		Action:    NewControllerAction(common),
		Dashboard: NewControllerDashboard(common),
		Map:       NewControllerMap(common),
	}
}

func (s *Controllers) Start() {
	s.Image.Start()
	s.Worker.Start()
	s.Action.Start()
	s.Dashboard.Start()
	s.Map.Start()
}

func (s *Controllers) Stop() {
	s.Image.Stop()
	s.Worker.Stop()
	s.Action.Stop()
	s.Dashboard.Stop()
	s.Map.Stop()
}
