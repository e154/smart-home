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

	"github.com/e154/smart-home/common"
)

// EntityPayload ...
type EntityPayload struct {
	AttributeSignature Attributes `json:"attribute_signature"`
}

// EntitySettings ...
type EntitySettings struct {
	Settings Attributes `json:"settings"`
}

// Entity ...
type Entity struct {
	Actions     []*EntityAction  `json:"actions"`
	States      []*EntityState   `json:"states"`
	Metrics     []*Metric        `json:"metrics"`
	Scripts     []*Script        `json:"scripts"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Id          common.EntityId  `json:"id"`
	Description string           `json:"description"`
	PluginName  string           `json:"plugin_name" validate:"required"`
	Icon        *string          `json:"icon"`
	Image       *Image           `json:"image"`
	ImageId     *int64           `json:"image_id"`
	Area        *Area            `json:"area"`
	AreaId      *int64           `json:"area_id"`
	Attributes  Attributes       `json:"attributes"`
	Settings    Attributes       `json:"settings"`
	ParentId    *common.EntityId `json:"parent_id"`
	Hidden      bool             `json:"hidden"`
	AutoLoad    bool             `json:"auto_load"`
	IsLoaded    bool             `json:"is_loaded"`
}
