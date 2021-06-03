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
type NewEntityState struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name" valid:"MaxSize(254);Required"`
	Description string  `json:"description"`
	Icon        *string `json:"icon"`
	Image       *Image  `json:"image"`
	ImageId     *int64  `json:"image_id" valid:"Required"`
	Style       string  `json:"style"`
}

// swagger:model
type EntityState struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" valid:"MaxSize(254);Required"`
	Description string    `json:"description"`
	Icon        *string   `json:"icon"`
	EntityId    string    `json:"entity_id" valid:"Required"`
	Image       *Image    `json:"image"`
	ImageId     *int64    `json:"image_id" valid:"Required"`
	Style       string    `json:"style"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// swagger:model
type EntityStatesShort struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Icon        *string `json:"icon"`
	Image       *Image  `json:"image"`
}
