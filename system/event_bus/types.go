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

package event_bus

import (
	"bytes"
	"encoding/json"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"time"
)

type EntityState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

type EventEntityState struct {
	EntityId    common.EntityId    `json:"entity_id"`
	Value       interface{}        `json:"value"`
	State       *EntityState       `json:"state"`
	Attributes  m.EntityAttributes `json:"attributes"`
	LastChanged *time.Time         `json:"last_changed"`
	LastUpdated *time.Time         `json:"last_updated"`
}

func (e1 EventEntityState) Compare(e2 EventEntityState) (ident bool) {

	if e1.State != e2.State {
		return
	}
	b1, _ := json.Marshal(e1.Attributes)
	b2, _ := json.Marshal(e2.Attributes)

	// Compare returns an integer comparing two byte slices lexicographically.
	// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
	// A nil argument is equivalent to an empty slice.
	if bytes.Equal(b1, b2) {
		ident = true
		return
	}

	return
}

type EventType string

const (
	TopicEntities = "entities"
)
