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

package dto

import (
	"github.com/e154/smart-home/internal/api/stub"
	m "github.com/e154/smart-home/pkg/models"
)

// Tag ...
type Tag struct{}

// NewTagDto ...
func NewTagDto() Tag {
	return Tag{}
}

// ToSearchResult ...
func (s Tag) ToSearchResult(list []*m.Tag) *stub.ApiSearchTagListResult {

	items := make([]stub.ApiTag, 0, len(list))

	for _, tag := range list {
		items = append(items, stub.ApiTag{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}

	return &stub.ApiSearchTagListResult{
		Items: items,
	}
}

// ToListResult ...
func (s Tag) ToListResult(list []*m.Tag) []*stub.ApiTag {

	items := make([]*stub.ApiTag, 0, len(list))

	for _, script := range list {
		items = append(items, &stub.ApiTag{
			Id:   script.Id,
			Name: script.Name,
		})
	}

	return items
}
