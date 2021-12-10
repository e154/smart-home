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

// MapLayer ...
type MapLayer struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name" validate:"max=254;required"`
	Description string        `json:"description"`
	Map         *Map          `json:"map"`
	MapId       int64         `json:"map_id" validate:"required"`
	Status      string        `json:"status" validate:"required"`
	Weight      int64         `json:"weight"`
	Elements    []*MapElement `json:"elements"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// SortMapLayersByWeight ...
type SortMapLayersByWeight []*MapLayer

// Len ...
func (l SortMapLayersByWeight) Len() int { return len(l) }

// Swap ...
func (l SortMapLayersByWeight) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Less ...
func (l SortMapLayersByWeight) Less(i, j int) bool { return l[i].Weight < l[j].Weight }

// SortMapLayer ...
type SortMapLayer struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}
