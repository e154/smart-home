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
	"encoding/json"
	"github.com/e154/smart-home/system/uuid"
	"github.com/e154/smart-home/system/validation"
	"time"
)

// Connection ...
type Connection struct {
	Uuid          uuid.UUID       `json:"uuid"`
	Name          string          `json:"name" valid:"MaxSize(254)"`
	ElementFrom   uuid.UUID       `json:"element_from" valid:"Required"`
	ElementTo     uuid.UUID       `json:"element_to" valid:"Required"`
	PointFrom     int64           `json:"point_from" valid:"Required"`
	PointTo       int64           `json:"point_to" valid:"Required"`
	Flow          *Flow           `json:"flow"`
	FlowId        int64           `json:"flow_id" valid:"Required"`
	Direction     string          `json:"direction"`
	GraphSettings json.RawMessage `json:"graph_settings"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// Valid ...
func (d *Connection) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
