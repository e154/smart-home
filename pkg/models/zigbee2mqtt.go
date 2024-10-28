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

package models

import (
	"time"
)

type Zigbee2mqttInfo struct {
	Networkmap    string     `json:"networkmap"`
	Status        string     `json:"status"`
	LastScan      *time.Time `json:"last_scan"`
	ScanInProcess bool       `json:"scan_in_process"`
}

// Zigbee2mqtt ...
type Zigbee2mqtt struct {
	Devices           []*Zigbee2mqttDevice `json:"devices"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	Name              string               `json:"name" validate:"required,max=254"`
	Login             string               `json:"login"`
	EncryptedPassword string               `json:"encrypted_password"`
	BaseTopic         string               `json:"base_topic"`
	Id                int64                `json:"id"`
	Password          *string              `json:"password"`
	Info              *Zigbee2mqttInfo     `json:"info"`
	PermitJoin        bool                 `json:"permit_join"`
}
