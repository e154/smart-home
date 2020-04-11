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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/graceful_service"
	metrics2 "github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

var (
	log = common.MustGetLogger("websocket.server")
)

// WebSocket ...
type WebSocket struct {
	Controllers *Controllers
}

// NewWebSocket ...
func NewWebSocket(adaptors *adaptors.Adaptors,
	stream *stream.StreamService,
	endpoint *endpoint.Endpoint,
	scripts *scripts.ScriptService,
	core *core.Core,
	graceful *graceful_service.GracefulService,
	metrics *metrics2.MetricManager,
	mqtt *mqtt.Mqtt,
	zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt) *WebSocket {

	server := &WebSocket{
		Controllers: NewControllers(adaptors, stream, scripts, core, endpoint, metrics, mqtt, zigbee2mqtt),
	}

	graceful.Subscribe(server)

	return server
}

// Start ...
func (s *WebSocket) Start() {
	log.Info("Serving websocket service")
	s.Controllers.Start()
}

// Shutdown ...
func (s *WebSocket) Shutdown() {
	log.Info("Server exiting")
	s.Controllers.Stop()
}
