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

package models

// DevModBusRtuConfig ...
type DevModBusRtuConfig struct {
	SlaveId  int    `json:"slave_id" mapstructure:"slave_id"`   // 1-32
	Baud     int    `json:"baud"`                               // 9600, 19200, ...
	DataBits int    `json:"data_bits" mapstructure:"data_bits"` // 5-9
	StopBits int    `json:"stop_bits" mapstructure:"stop_bits"` // 1, 2
	Parity   string `json:"parity"`                             // none, odd, even
	Timeout  int    `json:"timeout"`                            // milliseconds
}

// DevModBusTcpConfig ...
type DevModBusTcpConfig struct {
	SlaveId     int    `json:"slave_id"`
	AddressPort string `json:"address_port"`
}

// DevSmartBusConfig ...
type DevSmartBusConfig struct {
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required" mapstructure:"stop_bits"`
	Sleep    int `json:"sleep"`
}

// DevCommandConfig ...
type DevCommandConfig struct {
}

// DevMqttConfig ...
type DevMqttConfig struct {
	Address string `json:"address"`
}

// DevZigbee2mqttConfig ...
type DevZigbee2mqttConfig struct {
	Zigbee2mqttDeviceId string `json:"zigbee2mqtt_device_id"`
}

// An AllOfModel is composed out of embedded structs but it should build
// an allOf property
type DeviceProperties struct {
	// swagger:allOf
	*DevModBusRtuConfig
	// swagger:allOf
	*DevModBusTcpConfig
	// swagger:allOf
	*DevSmartBusConfig
	// swagger:allOf
	*DevCommandConfig
	// swagger:allOf
	*DevMqttConfig
	// swagger:allOf
	*DevZigbee2mqttConfig
}
