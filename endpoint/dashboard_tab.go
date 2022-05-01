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
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
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
func (t *DashboardTabEndpoint) Add(ctx context.Context, tab *m.DashboardTab) (result *m.DashboardTab, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = t.validation.Valid(tab); !ok {
		return
	}

	var id int64
	if id, err = t.adaptors.DashboardTab.Add(tab); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = t.adaptors.DashboardTab.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// GetById ...
func (t *DashboardTabEndpoint) GetById(ctx context.Context, id int64) (tab *m.DashboardTab, err error) {

	if tab, err = t.adaptors.DashboardTab.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	err = t.preloadEntities(tab)

	return
}

// Update ...
func (i *DashboardTabEndpoint) Update(ctx context.Context, params *m.DashboardTab) (result *m.DashboardTab, errs validator.ValidationErrorsTranslations, err error) {

	var board *m.DashboardTab
	if board, err = i.adaptors.DashboardTab.GetById(params.Id); err != nil {
		return
	}

	if err = copier.Copy(&board, &params); err != nil {
		return
	}

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	if err = i.adaptors.DashboardTab.Update(board); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = i.adaptors.DashboardTab.GetById(params.Id); err != nil {
		if !errors.Is(err, common.ErrNotFound) {
			err = errors.Wrap(common.ErrInternal, err.Error())
		}
	}

	return
}

// GetList ...
func (t *DashboardTabEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.DashboardTab, total int64, err error) {

	list, total, err = t.adaptors.DashboardTab.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	for _, tab := range list {
		err = t.preloadEntities(tab)
	}

	return
}

// Delete ...
func (t *DashboardTabEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = t.adaptors.DashboardTab.GetById(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	err = t.adaptors.DashboardTab.Delete(id)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

func (c *DashboardTabEndpoint) preloadEntities(tab *m.DashboardTab) (err error) {

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
	for entityId, _ := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entites []*m.Entity
	if entites, err = c.adaptors.Entity.GetByIds(entityIds); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	for _, entity := range entites {
		entityMap[entity.Id] = entity
	}

	tab.Entities = entityMap

	return
}
