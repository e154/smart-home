package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
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
}

func NewControllers(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	scripts *scripts.ScriptService,
	core *core.Core,
	endpoint *endpoint.Endpoint) *Controllers {
	common := NewControllerCommon(adaptors, stream, endpoint, scripts, core)
	return &Controllers{
		Image:     NewControllerImage(common),
		Worker:    NewControllerWorker(common),
		Action:    NewControllerAction(common),
		Dashboard: NewControllerDashboard(common),
	}
}
