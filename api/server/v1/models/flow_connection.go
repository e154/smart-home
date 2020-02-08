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
	"time"
)

// swagger:model
type FlowConnection struct {
	Uuid          string    `json:"uuid"`
	Name          string    `json:"name" valid:"MaxSize(254);Required"`
	ElementFrom   string    `json:"element_from" valid:"Required"`
	ElementTo     string    `json:"element_to" valid:"Required"`
	PointFrom     int64     `json:"point_from" valid:"Required"`
	PointTo       int64     `json:"point_to" valid:"Required"`
	FlowId        int64     `json:"flow_id" valid:"Required"`
	Direction     string    `json:"direction"`
	GraphSettings string    `json:"graph_settings"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
