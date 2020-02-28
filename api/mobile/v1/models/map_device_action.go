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

type MapDeviceAction struct {
	Id             int64         `json:"id"`
	DeviceAction   *DeviceAction `json:"device_action"`
	DeviceActionId int64         `json:"device_action_id" valid:"Required"`
	MapDeviceId    int64         `json:"map_device_id" valid:"Required"`
	Image          *Image        `json:"image"`
	ImageId        int64         `json:"image_id" valid:"Required"`
	Type           string        `json:"type"`
}

func (n MapDeviceAction) MarshalJSON() (b []byte, err error) {

	data := map[string]interface{}{
		"id":               n.Id,
		"name":             n.DeviceAction.Name,
		"description":      n.DeviceAction.Description,
		"device_action_id": n.DeviceActionId,
		"image":            n.Image,
		"map_device_id":    n.MapDeviceId,
	}

	if n.DeviceAction != nil {
		data["device_id"] = n.DeviceAction.DeviceId
	}

	b, err = json.Marshal(data)

	return
}