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

package gate

import (
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/api/gate/controllers"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/zigbee2mqtt"
)

var (
	log = common.MustGetLogger("gate")
)

type Gate struct {
	Controllers *Controllers
	graceful    *graceful_service.GracefulService
}

func NewGate(adaptors *adaptors.Adaptors,
	endpoint *endpoint.Endpoint,
	scripts *scripts.ScriptService,
	core *core.Core,
	graceful *graceful_service.GracefulService,
	gate *gate_client.GateClient,
	metric *metrics.MetricManager,
	mqtt *mqtt.Mqtt,
	zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt) *Gate {

	server := &Gate{
		Controllers: NewControllers(adaptors, scripts, core, endpoint, gate, metric, mqtt, zigbee2mqtt),
		graceful:    graceful,
	}

	return server
}

func (s *Gate) Start() {
	log.Info("Serving gate websocket service")
	s.graceful.Subscribe(s)

	s.Controllers.Start()
}

func (s *Gate) Shutdown() {
	log.Info("Server exiting")
	s.Controllers.Stop()
}
