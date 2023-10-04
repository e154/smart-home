// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package container

import (
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_authenticator"
)

// Mqtt ...
type Mqtt struct {
	authenticator mqtt_authenticator.MqttAuthenticator
}

// NewMqtt ...
func NewMqtt(authenticator mqtt_authenticator.MqttAuthenticator) mqtt.MqttServ {
	return &Mqtt{
		authenticator: authenticator,
	}
}

// Shutdown ...
func (m Mqtt) Shutdown() error {
	return nil
}

// Start ...
func (m Mqtt) Start() {}

// Publish ...
func (m Mqtt) Publish(topic string, payload []byte, qos uint8, retain bool) error {
	return nil
}

// NewClient ...
func (m Mqtt) NewClient(name string) mqtt.MqttCli {
	return NewMqttCli()
}

// RemoveClient ...
func (m Mqtt) RemoveClient(name string) {}

// Admin ...
func (m Mqtt) Admin() mqtt.Admin {
	return nil
}

// Authenticator ...
func (m Mqtt) Authenticator() mqtt_authenticator.MqttAuthenticator {
	return m.authenticator
}
