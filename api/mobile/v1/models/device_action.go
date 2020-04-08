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

// DeviceActionScript ...
type DeviceActionScript struct {
	Id int64 `json:"id"`
} 
// DeviceActionDevice ...
type DeviceActionDevice struct {
	Id int64 `json:"id"`
}

// swagger:model
type NewDeviceAction struct {
	Name        string              `json:"name" valid:"MaxSize(254);Required"`
	Description string              `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
}

// swagger:model
type UpdateDeviceAction struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name" valid:"MaxSize(254);Required"`
	Description string              `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
}

// swagger:model
type DeviceAction struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name" valid:"MaxSize(254);Required"`
	Description string              `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	DeviceId    int64               `json:"device_id"`
	Script      *DeviceActionScript `json:"script"`
}
