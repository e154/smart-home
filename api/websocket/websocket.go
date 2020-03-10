// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package websocket

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/api/websocket/controllers"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/graceful_service"
	metrics2 "github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/scripts"
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
	scripts *scripts.ScriptService,
	core *core.Core,
	graceful *graceful_service.GracefulService,
	metrics *metrics2.MetricManager) *WebSocket {

	server := &WebSocket{
		Controllers: NewControllers(adaptors, stream, scripts, core, endpoint, metrics),
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
