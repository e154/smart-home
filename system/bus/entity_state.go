// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package bus

import (
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// EntityState ...
type EntityState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

// EventEntityState ...
type EventEntityState struct {
	EntityId    common.EntityId `json:"entity_id"`
	Value       interface{}     `json:"value"`
	State       *EntityState    `json:"state"`
	Attributes  m.Attributes    `json:"attributes"`
	Settings    m.Attributes    `json:"settings"`
	LastChanged *time.Time      `json:"last_changed"`
	LastUpdated *time.Time      `json:"last_updated"`
}

// Compare ...
func (e1 EventEntityState) Compare(e2 EventEntityState) (ident bool) {

	if e1.State == nil || e2.State == nil {
		return
	}

	if e1.State.Name != e2.State.Name {
		return
	}

	for k1, v1 := range e1.Attributes {
		if !e2.Attributes[k1].Compare(v1) {
			return false
		}
	}

	ident = true

	return
}
