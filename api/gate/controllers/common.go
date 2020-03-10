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
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
)

type ControllerCommon struct {
	adaptors *adaptors.Adaptors
	endpoint *endpoint.Endpoint
	core     *core.Core
	scripts  *scripts.ScriptService
	gate     *gate_client.GateClient
	metric   *metrics.MetricManager
}

func NewControllerCommon(adaptors *adaptors.Adaptors,
	endpoint *endpoint.Endpoint,
	scripts *scripts.ScriptService,
	core *core.Core,
	gate *gate_client.GateClient,
	metric *metrics.MetricManager) *ControllerCommon {
	return &ControllerCommon{
		adaptors: adaptors,
		endpoint: endpoint,
		core:     core,
		scripts:  scripts,
		gate:     gate,
		metric:   metric,
	}
}

func (c *ControllerCommon) Err(client stream.IStreamClient, message stream.Message, err error) {
	msg := stream.Message{
		Id:      message.Id,
		Forward: stream.Response,
		Status:  stream.StatusError,
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	}
	client.Write(msg.Pack())
}
