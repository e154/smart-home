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
	"github.com/e154/smart-home/system/validation"
	"time"
)

type EntityPayload struct {
	AttributeSignature EntityAttributes `json:"attribute_signature"`
}

// Entity ...
type Entity struct {
	Id          common.EntityId   `json:"id"`
	Description string            `json:"description"`
	Type        common.EntityType `json:"type" valid:"Required"`
	Icon        *common.Icon      `json:"icon"`
	Image       *Image            `json:"image"`
	Actions     []*EntityAction   `json:"actions"`
	States      []*EntityState    `json:"states"`
	Area        *Area             `json:"area"`
	AreaId      *int64            `json:"area_id"`
	Metrics     []Metric          `json:"metrics"`
	Scripts     []Script          `json:"scripts"`
	Hidden      bool              `json:"hidden"`
	Attributes  EntityAttributes  `json:"attributes"`
	AutoLoad    bool              `json:"auto_load"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Valid ...
func (m *Entity) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}

// EntityPrototype ...
type IEntityPrototype interface {
	PrototypeName() common.EntityPrototypeType
}
