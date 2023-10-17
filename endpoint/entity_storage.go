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
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
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
func (i *EntityStorageEndpoint) GetList(ctx context.Context, entityId *common.EntityId, pagination common.PageParams, startDate, endDate *time.Time) (result *m.EntityStorageList, total int64, err error) {

	//var startDate, endDate *time.Time
	//if _startDate != nil {
	//	date, _ := time.Parse("2006-01-02", *_startDate)
	//	startDate = &date
	//}
	//if _endDate != nil {
	//	date, _ := time.Parse("2006-01-02", *_endDate)
	//	endDate = &date
	//}

	var items []*m.EntityStorage
	if items, total, err = i.adaptors.EntityStorage.ListByEntityId(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, entityId, startDate, endDate); err != nil {
		return
	}

	var idsMap = map[common.EntityId]struct{}{}

	for j := range items {
		idsMap[items[j].EntityId] = struct{}{}
	}

	var ids []common.EntityId
	ids = make([]common.EntityId, 0, len(idsMap))
	for id := range idsMap {
		ids = append(ids, id)
	}

	var entities []*m.Entity
	if entities, err = i.adaptors.Entity.GetByIdsSimple(ctx, ids); err != nil {
		return
	}

	attributes := make(map[common.EntityId]m.Attributes)

	for _, entity := range entities {
		attributes[entity.Id] = entity.Attributes
	}

	result = &m.EntityStorageList{
		Items:      items,
		Attributes: attributes,
	}

	return
}
