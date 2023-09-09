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
)

// DashboardCardItemEndpoint ...
type DashboardCardItemEndpoint struct {
	*CommonEndpoint
}

// NewDashboardCardItemEndpoint ...
func NewDashboardCardItemEndpoint(common *CommonEndpoint) *DashboardCardItemEndpoint {
	return &DashboardCardItemEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (c *DashboardCardItemEndpoint) Add(ctx context.Context, card *m.DashboardCardItem) (result *m.DashboardCardItem, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = c.validation.Valid(card); !ok {
		return
	}

	var id int64
	if id, err = c.adaptors.DashboardCardItem.Add(ctx, card); err != nil {
		return
	}

	result, err = c.adaptors.DashboardCardItem.GetById(ctx, id)

	return
}

// GetById ...
func (c *DashboardCardItemEndpoint) GetById(ctx context.Context, id int64) (card *m.DashboardCardItem, err error) {

	card, err = c.adaptors.DashboardCardItem.GetById(ctx, id)

	return
}

// Update ...
func (i *DashboardCardItemEndpoint) Update(ctx context.Context, params *m.DashboardCardItem) (result *m.DashboardCardItem, errs validator.ValidationErrorsTranslations, err error) {

	var board *m.DashboardCardItem
	if board, err = i.adaptors.DashboardCardItem.GetById(ctx, params.Id); err != nil {
		return
	}

	if err = copier.Copy(&board, &params); err != nil {
		return
	}

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	if err = i.adaptors.DashboardCardItem.Update(ctx, board); err != nil {
		return
	}

	result, err = i.adaptors.DashboardCardItem.GetById(ctx, params.Id)

	return
}

// GetList ...
func (c *DashboardCardItemEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.DashboardCardItem, total int64, err error) {

	list, total, err = c.adaptors.DashboardCardItem.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)

	return
}

// Delete ...
func (c *DashboardCardItemEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = c.adaptors.DashboardCardItem.GetById(ctx, id)
	if err != nil {
		return
	}

	err = c.adaptors.DashboardCardItem.Delete(ctx, id)

	return
}
