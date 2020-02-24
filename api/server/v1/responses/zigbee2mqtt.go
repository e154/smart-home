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

package responses

import (
	"github.com/e154/smart-home/api/server/v1/models"
)

// swagger:response Zigbee2mqttList
type Zigbee2mqttList struct {
	// in:body
	Body struct {
		Items []*models.Zigbee2mqtt `json:"items"`
		Meta  struct {
			Limit       int64 `json:"limit"`
			ObjectCount int64 `json:"object_count"`
			Offset      int64 `json:"offset"`
		} `json:"meta"`
	}
}

// swagger:response Zigbee2mqttDeviceSearch
type Zigbee2mqttDeviceSearch struct {
	// in:body
	Body struct {
		Zigbee2mqttDevices []*models.Zigbee2mqttDeviceShort `json:"devices"`
	}
}
