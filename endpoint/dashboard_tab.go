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
	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
)

// DashboardTabEndpoint ...
type DashboardTabEndpoint struct {
	*CommonEndpoint
}

// NewDashboardTabEndpoint ...
func NewDashboardTabEndpoint(common *CommonEndpoint) *DashboardTabEndpoint {
	return &DashboardTabEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (t *DashboardTabEndpoint) Add(ctx context.Context, tab *m.DashboardTab) (result *m.DashboardTab, err error) {

	if ok, errs := t.validation.Valid(tab); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var id int64
	if id, err = t.adaptors.DashboardTab.Add(ctx, tab); err != nil {
		return
	}

	result, err = t.adaptors.DashboardTab.GetById(ctx, id)

	return
}

// GetById ...
func (t *DashboardTabEndpoint) GetById(ctx context.Context, id int64) (tab *m.DashboardTab, err error) {

	if tab, err = t.adaptors.DashboardTab.GetById(ctx, id); err != nil {
		return
	}

	err = t.preloadEntities(ctx, tab)

	return
}

// Update ...
func (t *DashboardTabEndpoint) Update(ctx context.Context, params *m.DashboardTab) (result *m.DashboardTab, err error) {

	var tab *m.DashboardTab
	if tab, err = t.adaptors.DashboardTab.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&tab, &params); err != nil {
		return
	}

	if ok, errs := t.validation.Valid(tab); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = t.adaptors.DashboardTab.Update(ctx, tab); err != nil {
		return
	}

	result, err = t.adaptors.DashboardTab.GetById(ctx, params.Id)

	return
}

// GetList ...
func (t *DashboardTabEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.DashboardTab, total int64, err error) {

	list, total, err = t.adaptors.DashboardTab.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		return
	}

	for _, tab := range list {
		err = t.preloadEntities(ctx, tab)
	}

	return
}

// Delete ...
func (t *DashboardTabEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = t.adaptors.DashboardTab.GetById(ctx, id)
	if err != nil {
		return
	}

	err = t.adaptors.DashboardTab.Delete(ctx, id)

	return
}

func (t *DashboardTabEndpoint) preloadEntities(ctx context.Context, tab *m.DashboardTab) (err error) {

	// get child entities
	entityMap := make(map[common.EntityId]*m.Entity)
	for _, card := range tab.Cards {
		for _, item := range card.Items {
			if item.EntityId != nil {
				entityMap[*item.EntityId] = nil
			}
		}
	}

	entityIds := make([]common.EntityId, 0, len(entityMap))
	for entityId := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entities []*m.Entity
	if entities, err = t.adaptors.Entity.GetByIds(ctx, entityIds); err != nil {
		return
	}

	for _, entity := range entities {
		entityMap[entity.Id] = entity
	}

	tab.Entities = entityMap

	return
}
