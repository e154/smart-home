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

// MapDevice ...
type MapDevice struct {
	Id         int64              `json:"id"`
	SystemName string             `json:"system_name" valid:"Required"`
	Device     *Device            `json:"device"`
	DeviceId   int64              `json:"device_id" valid:"Required"`
	Image      *Image             `json:"image"`
	ImageId    int64              `json:"image_id"`
	Actions    []*MapDeviceAction `json:"actions"`
	States     []*MapDeviceState  `json:"states"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
