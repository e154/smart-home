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
	"encoding/json"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/null"
	"github.com/e154/smart-home/system/validation"
	"time"
)

// MapElementGraphSettingsPosition ...
type MapElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}

// MapElementGraphSettings ...
type MapElementGraphSettings struct {
	Width    null.Int64                      `json:"width"`
	Height   null.Int64                      `json:"height"`
	Position MapElementGraphSettingsPosition `json:"position"`
}

// MapElementPrototype ...
type MapElementPrototype struct {
	*MapImage
	*MapText
	*Entity
}

// MarshalJSON ...
func (n MapElementPrototype) MarshalJSON() (b []byte, err error) {

	switch {
	case n.MapText != nil:
		b, err = json.Marshal(n.MapText)
	case n.MapImage != nil:
		b, err = json.Marshal(n.MapImage)
	case n.Entity != nil:
		b, err = json.Marshal(n.Entity)
	default:
		b, err = json.Marshal(struct{}{})
		return
	}
	return
}

// UnmarshalJSON ...
func (n *MapElementPrototype) UnmarshalJSON(data []byte) (err error) {

	entity := &Entity{}
	err = json.Unmarshal(data, entity)
	if ok, _ := entity.Valid(); ok {
		n.Entity = entity
		return
	}

	image := &MapImage{}
	err = json.Unmarshal(data, image)
	if image.ImageId != 0 {
		n.MapImage = image
		return
	}

	text := &MapText{}
	err = json.Unmarshal(data, text)
	n.MapText = text
	return
}

// Entity ...
type MapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name" valid:"Required"`
	Description   string                  `json:"description"`
	PrototypeId   interface{}             `json:"prototype_id"`
	PrototypeType MapElementPrototypeType `json:"prototype_type"`
	Prototype     MapElementPrototype     `json:"prototype" valid:"Required"`
	MapId         int64                   `json:"map_id" valid:"Required"`
	LayerId       int64                   `json:"layer_id" valid:"Required"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        StatusType              `json:"status" valid:"Required"`
	Weight        int64                   `json:"weight"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

// Valid ...
func (m *MapElement) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}

// SortMapElementByWeight ...
type SortMapElementByWeight []*MapElement

// Len ...
func (l SortMapElementByWeight) Len() int { return len(l) }

// Swap ...
func (l SortMapElementByWeight) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Less ...
func (l SortMapElementByWeight) Less(i, j int) bool { return l[i].Weight < l[j].Weight }

// SortMapElement ...
type SortMapElement struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}
