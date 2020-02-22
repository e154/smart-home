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

package zigbee2mqtt

import (
	"github.com/e154/smart-home/adaptors"
	models2 "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("zigbee2mqtt")
)

type Zigbee2mqtt struct {
	graceful  *graceful_service.GracefulService
	mqtt      *mqtt.Mqtt
	adaptors  *adaptors.Adaptors
	bridges   []*Bridge
	isStarted bool
}

func NewZigbee2mqtt(graceful *graceful_service.GracefulService,
	mqtt *mqtt.Mqtt,
	adaptors *adaptors.Adaptors) *Zigbee2mqtt {
	return &Zigbee2mqtt{
		graceful: graceful,
		mqtt:     mqtt,
		bridges:  make([]*Bridge, 0),
		adaptors: adaptors,
	}
}

func (z *Zigbee2mqtt) Start() {
	if z.isStarted {
		return
	}
	z.isStarted = true

	models, _, err := z.adaptors.Zigbee2mqtt.List(99, 0)
	if err != nil {
		log.Error(err.Error())
	}

	if len(models) == 0 {
		model := &models2.Zigbee2mqtt{
			Name:       "zigbee2mqtt",
			BaseTopic:  "zigbee2mqtt",
			PermitJoin: true,
		}
		model.Id, err = z.adaptors.Zigbee2mqtt.Add(model)
		if err != nil {
			log.Error(err.Error())
			return
		}
		models = append(models, model)
	}

	for _, model := range models {
		bridge := NewBridge(z.mqtt, z.adaptors, model)
		bridge.Start()
		z.bridges = append(z.bridges, bridge)
	}
}

func (z *Zigbee2mqtt) Shutdown() {
	if !z.isStarted {
		return
	}
	z.isStarted = false
	for _, bridge := range z.bridges {
		bridge.Stop()
	}
}
