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
	"fmt"
	"time"

	"github.com/e154/smart-home/internal/common"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/e154/smart-home/pkg/models"
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
func (i *EntityStorageEndpoint) GetList(ctx context.Context, entityIds []pkgCommon.EntityId, pagination common.PageParams,
	startDate, endDate *time.Time) (result *models.EntityStorageList, total int64, err error) {

	keyResult := fmt.Sprintf("entity_storage_%d_%d_%s_%s_%v_%s_%s", pagination.Limit, pagination.Offset,
		pagination.Order, pagination.SortBy, entityIds, startDate, endDate)

	keyTotal := fmt.Sprintf("%s_total", keyResult)

	if ok, _ := i.cache.IsExist(ctx, keyResult); ok {
		v, _ := i.cache.Get(ctx, keyResult)
		result = v.(*models.EntityStorageList)
		vTotal, _ := i.cache.Get(ctx, keyTotal)
		total = vTotal.(int64)
		return
	}

	var items []*models.EntityStorage
	if items, total, err = i.adaptors.EntityStorage.List(ctx, pagination.Limit, pagination.Offset,
		pagination.Order, pagination.SortBy, entityIds, startDate, endDate); err != nil {
		return
	}

	var entitiesList = map[pkgCommon.EntityId]*models.Entity{}

	for j := range items {
		entitiesList[items[j].EntityId] = nil
	}

	if len(entityIds) == 0 {
		for id := range entitiesList {
			entityIds = append(entityIds, id)
		}
	}

	var entities []*models.Entity
	if entities, err = i.adaptors.Entity.GetByIdsSimple(ctx, entityIds); err != nil {
		return
	}

	attributes := make(map[pkgCommon.EntityId]models.Attributes)

	for _, entity := range entities {
		entitiesList[entity.Id] = entity
		attributes[entity.Id] = entity.Attributes
	}

	var ok bool
	var entity *models.Entity
	for _, item := range items {
		item.StateDescription = item.State
		if entity, ok = entitiesList[item.EntityId]; !ok {
			continue
		}
		item.EntityDescription = entity.Id.String()
		if entity.Description != "" {
			item.EntityDescription = entity.Description
		}
		for _, state := range entity.States {
			if state.Name == item.State && state.Description != "" {
				item.StateDescription = state.Description
				continue
			}
		}
	}

	result = &models.EntityStorageList{
		Items:      items,
		Attributes: attributes,
	}

	_ = i.cache.Put(ctx, keyResult, result, 10*time.Second)
	_ = i.cache.Put(ctx, keyTotal, total, 10*time.Second)

	return
}
