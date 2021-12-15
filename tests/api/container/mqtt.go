// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

type Mqtt struct {
	authenticator mqtt_authenticator.MqttAuthenticator
}

func NewMqtt(authenticator mqtt_authenticator.MqttAuthenticator) mqtt.MqttServ {
	return &Mqtt{
		authenticator: authenticator,
	}
}

func (m Mqtt) Shutdown() error {
	return nil
}

func (m Mqtt) Start() {}

func (m Mqtt) Publish(topic string, payload []byte, qos uint8, retain bool) error {
	return nil
}

func (m Mqtt) NewClient(name string) mqtt.MqttCli {
	return NewMqttCli()
}

func (m Mqtt) RemoveClient(name string) {}

func (m Mqtt) Admin() mqtt.Admin {
	return nil
}

func (m Mqtt) Authenticator() mqtt_authenticator.MqttAuthenticator {
	return m.authenticator
}
