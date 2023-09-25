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

package mqtt_client

// Config ...
type Config struct {
	KeepAlive      int    `json:"keep_alive"`
	PingTimeout    int    `json:"ping_timeout"`
	Broker         string `json:"broker"`
	ClientID       string `json:"client_id"`
	ConnectTimeout int    `json:"connect_timeout"`
	CleanSession   bool   `json:"clean_session"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Qos            byte   `json:"qos"`
}
