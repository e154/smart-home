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

package endpoint

import (
	"context"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/pkg/errors"
	"time"
)

// EntityStorageEndpoint ...
type EntityStorageEndpoint struct {
	*CommonEndpoint
}

// NewEntityStorageEndpoint ...
func NewEntityStorageEndpoint(common *CommonEndpoint) *EntityStorageEndpoint {
	return &EntityStorageEndpoint{
		CommonEndpoint: common,
	}
}

// GetList ...
func (i *EntityStorageEndpoint) GetList(ctx context.Context, entityId common.EntityId, pagination common.PageParams, _startDate, _endDate *string) (items []*m.EntityStorage, total int64, err error) {

	var startDate, endDate *time.Time
	if _startDate != nil {
		date, _ := time.Parse("2006-01-02", *_startDate)
		startDate = &date
	}
	if _endDate != nil {
		date, _ := time.Parse("2006-01-02", *_endDate)
		endDate = &date
	}

	items, total, err = i.adaptors.EntityStorage.ListByEntityId(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, entityId, startDate, endDate)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}
