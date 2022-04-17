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

// DashboardEndpoint ...
type DashboardEndpoint struct {
	*CommonEndpoint
}

// NewDashboardEndpoint ...
func NewDashboardEndpoint(common *CommonEndpoint) *DashboardEndpoint {
	return &DashboardEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (d *DashboardEndpoint) Add(ctx context.Context, board *m.Dashboard) (result *m.Dashboard, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = d.validation.Valid(board); !ok {
		return
	}

	var id int64
	if id, err = d.adaptors.Dashboard.Add(board); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = d.adaptors.Dashboard.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// GetById ...
func (d *DashboardEndpoint) GetById(ctx context.Context, id int64) (board *m.Dashboard, err error) {

	if board, err = d.adaptors.Dashboard.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	err = d.preloadEntities(board)

	return
}

// Update ...
func (i *DashboardEndpoint) Update(ctx context.Context, params *m.Dashboard) (result *m.Dashboard, errs validator.ValidationErrorsTranslations, err error) {

	var board *m.Dashboard
	if board, err = i.adaptors.Dashboard.GetById(params.Id); err != nil {
		return
	}

	if err = copier.Copy(&board, &params); err != nil {
		return
	}

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	if err = i.adaptors.Dashboard.Update(board); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = i.adaptors.Dashboard.GetById(params.Id); err != nil {
		if !errors.Is(err, common.ErrNotFound) {
			err = errors.Wrap(common.ErrInternal, err.Error())
		}
	}

	return
}

// GetList ...
func (d *DashboardEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.Dashboard, total int64, err error) {

	list, total, err = d.adaptors.Dashboard.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	for _, board := range list {
		err = d.preloadEntities(board)
	}

	return
}

// Delete ...
func (d *DashboardEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = d.adaptors.Dashboard.GetById(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	err = d.adaptors.Dashboard.Delete(id)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

func (c *DashboardEndpoint) preloadEntities(board *m.Dashboard) (err error) {

	// get child entities
	entityMap := make(map[common.EntityId]*m.Entity)
	for _, tab := range board.Tabs {
		for _, card := range tab.Cards {
			for _, item := range card.Items {
				entityMap[item.EntityId] = nil
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

	board.Entities = entityMap

	return
}
