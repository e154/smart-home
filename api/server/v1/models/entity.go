// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
type NewEntity struct {
	Id          string             `json:"id"`
	Description string             `json:"description"`
	Type        string             `json:"type" valid:"Required"`
	Icon        *string            `json:"icon"`
	Image       *Image             `json:"image"`
	Actions     []*NewEntityAction `json:"actions"`
	States      []*NewEntityState  `json:"states"`
	AreaId      *int64             `json:"area_id"`
	Metrics     []Metric           `json:"metrics"`
	Scripts     []Script           `json:"scripts"`
	Hidden      bool               `json:"hidden"`
	Attributes  Attributes         `json:"attributes"`
	Settings    Attributes         `json:"settings"`
	AutoLoad    bool               `json:"auto_load"`
}

// swagger:model
type UpdateEntity struct {
	Description string          `json:"description"`
	Type        string          `json:"type" valid:"Required"`
	Icon        *string         `json:"icon"`
	Image       *Image          `json:"image"`
	Actions     []*EntityAction `json:"actions"`
	States      []*EntityState  `json:"states"`
	AreaId      *int64          `json:"area_id"`
	Metrics     []Metric        `json:"metrics"`
	Scripts     []Script        `json:"scripts"`
	Hidden      bool            `json:"hidden"`
	Attributes  Attributes      `json:"attributes"`
	Settings    Attributes      `json:"settings"`
	AutoLoad    bool            `json:"auto_load"`
}

// swagger:model
type Entity struct {
	Id          string          `json:"id"`
	Description string          `json:"description"`
	Type        string          `json:"type" valid:"Required"`
	Icon        *string         `json:"icon"`
	Image       *Image          `json:"image"`
	Actions     []*EntityAction `json:"actions"`
	States      []*EntityState  `json:"states"`
	AreaId      *int64          `json:"area_id"`
	Metrics     []Metric        `json:"metrics"`
	Scripts     []Script        `json:"scripts"`
	Hidden      bool            `json:"hidden"`
	Attributes  Attributes      `json:"attributes"`
	Settings    Attributes      `json:"settings"`
	AutoLoad    bool            `json:"auto_load"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
