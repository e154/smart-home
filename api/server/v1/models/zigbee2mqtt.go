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

import "time"

// swagger:model
type NewZigbee2mqtt struct {
	Name           string `json:"name"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	PermitJoin     bool   `json:"permit_join"`
	BaseTopic      string `json:"base_topic"`
}

// swagger:model
type UpdateZigbee2mqtt struct {
	Name           string `json:"name"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
	PermitJoin     bool   `json:"permit_join"`
	BaseTopic      string `json:"base_topic"`
}

// swagger:model
type Zigbee2mqtt struct {
	Id         int64                `json:"id"`
	Name       string               `json:"name"`
	Login      string               `json:"login"`
	Devices    []*Zigbee2mqttDevice `json:"devices"`
	PermitJoin bool                 `json:"permit_join"`
	BaseTopic  string               `json:"base_topic"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
}

// swagger:model
type Zigbee2mqttInfo struct {
	ScanInProcess bool        `json:"scan_in_process"`
	LastScan      time.Time   `json:"last_scan"`
	Networkmap    string      `json:"networkmap"`
	Status        string      `json:"status"`
	Model         Zigbee2mqtt `json:"model"`
}
