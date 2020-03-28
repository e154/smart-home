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

package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeZigbee2mqtt = DeviceType("zigbee2mqtt")
)

type DevZigbee2mqttConfig struct {
	Validation
	Zigbee2mqttDeviceId string `json:"zigbee2mqtt_device_id"`
}

type DevZigbee2mqttRequest struct {
	Path    string `json:"path"`
	Payload []byte `json:"payload"`
}

// params:
// result
// error
// time
type DevZigbee2mqttResponse struct {
	BaseResponse
}

// Javascript Binding
//
// Zigbee2mqtt(path, payload)
//
func NewZigbee2mqttBind(path string, payload string) Zigbee2mqttBind {
	return Zigbee2mqttBind{
		Path:    path,
		Payload: []byte(payload),
	}
}

type Zigbee2mqttBind struct {
	Path    string
	Payload []byte
}
