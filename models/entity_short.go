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
	"github.com/e154/smart-home/common"
)

// todo: check
// EntityShort ...
type EntityShort struct {
	Id          common.EntityId     `json:"id"`
	Type        string              `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Icon        *string             `json:"icon"`
	ImageUrl    *string             `json:"image_url"`
	Actions     []EntityActionShort `json:"actions"`
	States      []EntityStateShort  `json:"states"`
	State       *EntityStateShort   `json:"state"`
	Attributes  Attributes          `json:"attributes"`
	Settings    Attributes          `json:"settings"`
	Area        *Area               `json:"area"`
	Metrics     []*Metric           `json:"metrics"`
	Hidden      bool                `json:"hidden"`
}
