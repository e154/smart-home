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

import "encoding/json"

type MapDeviceState struct {
	Id            int64        `json:"id"`
	DeviceState   *DeviceState `json:"device_state"`
	DeviceStateId int64        `json:"device_state_id" valid:"Required"`
	MapDeviceId   int64        `json:"map_device_id" valid:"Required"`
	Image         *Image       `json:"image"`
	ImageId       int64        `json:"image_id" valid:"Required"`
	Style         string       `json:"style"`
}

func (n MapDeviceState) MarshalJSON() (b []byte, err error) {

	data := map[string]interface{}{
		"id":              n.Id,
		"system_name":     n.DeviceState.SystemName,
		"description":     n.DeviceState.Description,
		"device_state_id": n.DeviceStateId,
		"map_device_id":   n.MapDeviceId,
		"image":           n.Image,
		"style":           n.Style,
	}

	if n.DeviceState != nil {
		data["device_id"] = n.DeviceState.DeviceId
	}

	b, err = json.Marshal(data)

	return
}
