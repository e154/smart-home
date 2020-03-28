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

import (
	"github.com/e154/smart-home/common"
	"time"
)

type MapDeviceHistory struct {
	Id           int64                `json:"id"`
	MapDeviceId  int64                `json:"map_device_id"`
	MapElement   *MapElement          `json:"map_element"`
	MapElementId int64                `json:"map_element_id"`
	LogLevel     common.LogLevel      `json:"log_level"`
	Type         common.MapDeviceHistoryType `json:"type"`
	Description  string               `json:"description"`
	CreatedAt    time.Time            `json:"created_at"`
}
