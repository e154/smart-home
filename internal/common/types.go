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

package common

// PageParams ...
type PageParams struct {
	Limit   int64  `json:"limit" validate:"required,gte=1,lte=1000"`
	Offset  int64  `json:"offset" validate:"required,gte=0,lte=1000"`
	Order   string `json:"order" validate:"required,oneof=created_at"`
	SortBy  string `json:"sort_by" validate:"required,oneof=desc asc"`
	PageReq int64
	SortReq string
}

// SearchParams ...
type SearchParams struct {
	Query  string `json:"query" validate:"required,min=1,max;255"`
	Limit  int64  `json:"limit" validate:"required,gte=1,lte=1000"`
	Offset int64  `json:"offset" validate:"required,gte=0,lte=1000"`
}
