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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// TagEndpoint ...
type TagEndpoint struct {
	*CommonEndpoint
}

// NewTagEndpoint ...
func NewTagEndpoint(common *CommonEndpoint) *TagEndpoint {
	return &TagEndpoint{
		CommonEndpoint: common,
	}
}

// GetList ...
func (n *TagEndpoint) GetList(ctx context.Context, pagination common.PageParams, query *string, names *[]string) (result []*m.Tag, total int64, err error) {
	result, total, err = n.adaptors.Tag.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, query, names)
	return
}

// Search ...
func (n *TagEndpoint) Search(ctx context.Context, query string, limit, offset int64) (tags []*m.Tag, total int64, err error) {

	tags, total, err = n.adaptors.Tag.Search(ctx, query, limit, offset)

	return
}
