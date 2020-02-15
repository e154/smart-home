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
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_client"
	"github.com/op/go-logging"
	"sync"
)

var (
	log = logging.MustGetLogger("zigbee2mqtt")
)


type Zigbee2mqtt struct {
	graceful   *graceful_service.GracefulService
	mqtt       *mqtt.Mqtt
	mqttLock sync.Mutex
	mqttClient *mqtt_client.Client
}

func NewZigbee2mqtt(graceful *graceful_service.GracefulService,
	mqtt *mqtt.Mqtt) *Zigbee2mqtt {
	return &Zigbee2mqtt{
		graceful: graceful,
		mqtt:     mqtt,
	}
}

func (z *Zigbee2mqtt) Start() {
	z.graceful.Subscribe(z)

	z.mqttLock.Lock()
	defer z.mqttLock.Unlock()

	var err error
	if z.mqttClient == nil {
		log.Info("create new mqtt client...")
		if z.mqttClient, err = z.mqtt.NewClient(nil); err != nil {
			log.Error(err.Error())
		}
	}

	if !z.mqttClient.IsConnected() {
		if err = z.mqttClient.Connect(); err != nil {
			log.Error(err.Error())
		}
	}

	//// /home/node/resp
	//if err := z.mqttClient.Subscribe(z.topic("resp"), 0, z.onPublish); err != nil {
	//	log.Warning(err.Error())
	//}
	//
	//// /home/node/ping
	//if err := z.mqttClient.Subscribe(z.topic("ping"), 0, z.ping); err != nil {
	//	log.Warning(err.Error())
	//}

	log.Info("Starting...")
}

func (z *Zigbee2mqtt) Shutdown() {

	log.Info("Exiting...")
}
