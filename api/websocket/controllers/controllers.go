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

package controllers

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	metrics2 "github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = common.MustGetLogger("stream.controllers")
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
	metrics *metrics2.MetricManager) *Controllers {
	common := NewControllerCommon(adaptors, stream, endpoint, scripts, core, metrics)
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
