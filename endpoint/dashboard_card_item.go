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
	if id, err = c.adaptors.DashboardCardItem.Add(card); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = c.adaptors.DashboardCardItem.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// GetById ...
func (c *DashboardCardItemEndpoint) GetById(ctx context.Context, id int64) (card *m.DashboardCardItem, err error) {

	if card, err = c.adaptors.DashboardCardItem.GetById(id); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}

	return
}

// Update ...
func (i *DashboardCardItemEndpoint) Update(ctx context.Context, params *m.DashboardCardItem) (result *m.DashboardCardItem, errs validator.ValidationErrorsTranslations, err error) {

	var board *m.DashboardCardItem
	if board, err = i.adaptors.DashboardCardItem.GetById(params.Id); err != nil {
		return
	}

	if err = copier.Copy(&board, &params); err != nil {
		return
	}

	var ok bool
	if ok, errs = i.validation.Valid(params); !ok {
		return
	}

	if err = i.adaptors.DashboardCardItem.Update(board); err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	if result, err = i.adaptors.DashboardCardItem.GetById(params.Id); err != nil {
		if !errors.Is(err, common.ErrNotFound) {
			err = errors.Wrap(common.ErrInternal, err.Error())
		}
	}

	return
}

// GetList ...
func (c *DashboardCardItemEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []*m.DashboardCardItem, total int64, err error) {

	list, total, err = c.adaptors.DashboardCardItem.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}

// Delete ...
func (c *DashboardCardItemEndpoint) Delete(ctx context.Context, id int64) (err error) {

	_, err = c.adaptors.DashboardCardItem.GetById(id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return
		}
		err = errors.Wrap(common.ErrInternal, err.Error())
		return
	}

	err = c.adaptors.DashboardCardItem.Delete(id)
	if err != nil {
		err = errors.Wrap(common.ErrInternal, err.Error())
	}
	return
}
