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

// Zigbee2mqtt ...
type Zigbee2mqtt struct {
	Id                int64                `json:"id"`
	Name              string               `json:"name" validate:"required,max=254"`
	Login             string               `json:"login"`
	Password          *string              `json:"password"`
	EncryptedPassword string               `json:"encrypted_password"`
	Devices           []*Zigbee2mqttDevice `json:"devices"`
	PermitJoin        bool                 `json:"permit_join"`
	BaseTopic         string               `json:"base_topic"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}
