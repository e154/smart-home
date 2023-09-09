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

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	m "github.com/e154/smart-home/models"
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
	if id, err = d.adaptors.Dashboard.Add(ctx, board); err != nil {
		return
	}

	result, err = d.adaptors.Dashboard.GetById(ctx, id)

	return
}

// GetById ...
func (d *DashboardEndpoint) GetById(ctx context.Context, id int64) (board *m.Dashboard, err error) {

	if board, err = d.adaptors.Dashboard.GetById(ctx, id); err != nil {
		return
	}

	err = d.preloadEntities(ctx, board)

	return
}

// Search ...
func (d *DashboardEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Dashboard, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = d.adaptors.Dashboard.Search(ctx, query, limit, offset)

	return
}

// Update ...
func (i *DashboardEndpoint) Update(ctx context.Context, params *m.Dashboard) (result *m.Dashboard, errs validator.ValidationErrorsTranslations, err error) {

	var board *m.Dashboard
	if board, err = i.adaptors.Dashboard.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&board, &params); err != nil {
		return
	}

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	if err = i.adaptors.Dashboard.Update(ctx, board); err != nil {
		return
	}

	if result, err = i.adaptors.Dashboard.GetById(ctx, params.Id); err != nil {
		if !errors.Is(err, apperr.ErrNotFound) {
		}
	}

	return
}

// GetList ...
func (d *DashboardEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.Dashboard, total int64, err error) {

	list, total, err = d.adaptors.Dashboard.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		return
	}

	for _, board := range list {
		err = d.preloadEntities(ctx, board)
	}

	return
}

// Delete ...
func (d *DashboardEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = d.adaptors.Dashboard.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		return
	}

	err = d.adaptors.Dashboard.Delete(ctx, id)
	if err != nil {
	}
	return
}

func (d *DashboardEndpoint) preloadEntities(ctx context.Context, board *m.Dashboard) (err error) {

	// get child entities
	entityMap := make(map[common.EntityId]*m.Entity)
	for _, tab := range board.Tabs {
		for _, card := range tab.Cards {
			if card.EntityId != nil {
				entityMap[*card.EntityId] = nil
			}
			for _, item := range card.Items {
				if item.EntityId != nil {
					entityMap[*item.EntityId] = nil
				}
			}
		}
	}

	entityIds := make([]common.EntityId, 0, len(entityMap))
	for entityId := range entityMap {
		entityIds = append(entityIds, entityId)
	}

	var entities []*m.Entity
	if entities, err = d.adaptors.Entity.GetByIds(ctx, entityIds); err != nil {
		return
	}

	for _, entity := range entities {
		entityMap[entity.Id] = entity
	}

	board.Entities = entityMap

	return
}

// Import ...
func (d *DashboardEndpoint) Import(ctx context.Context, board *m.Dashboard) (result *m.Dashboard, err error) {

	//b, _ := json.Marshal(board)
	//fmt.Println(string(b))

	var id int64
	if id, err = d.adaptors.Dashboard.Import(ctx, board); err != nil {
		return
	}

	if result, err = d.adaptors.Dashboard.GetById(ctx, id); err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			return
		}
		return
	}

	err = d.preloadEntities(ctx, board)

	return
}
