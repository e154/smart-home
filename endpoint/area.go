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

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/go-playground/validator/v10"
)

// AreaEndpoint ...
type AreaEndpoint struct {
	*CommonEndpoint
}

// NewAreaEndpoint ...
func NewAreaEndpoint(common *CommonEndpoint) *AreaEndpoint {
	return &AreaEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *AreaEndpoint) Add(ctx context.Context, params *m.Area) (result *m.Area, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetContext(err, errs)
		return
	}

	if _, err = n.adaptors.Area.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Area.GetByName(params.Name)

	return
}

// GetById ...
func (n *AreaEndpoint) GetById(ctx context.Context, id int64) (result *m.Area, err error) {

	result, err = n.adaptors.Area.GetById(id)

	return
}

// GetByName ...
func (n *AreaEndpoint) GetByName(ctx context.Context, name string) (result *m.Area, err error) {

	result, err = n.adaptors.Area.GetByName(name)

	return
}

// Update ...
func (n *AreaEndpoint) Update(ctx context.Context, params *m.Area) (result *m.Area, errs validator.ValidationErrorsTranslations, err error) {

	var area *m.Area
	area, err = n.adaptors.Area.GetById(params.Id)
	if err != nil {
		return
	}

	area.Name = params.Name
	area.Description = params.Description

	var ok bool
	if ok, errs = n.validation.Valid(area); !ok {
		return
	}

	if err = n.adaptors.Area.Update(area); err != nil {
		return
	}

	result = &m.Area{
		Id:          area.Id,
		Name:        area.Name,
		Description: area.Description,
	}

	return
}

// GetList ...
func (n *AreaEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Area, total int64, err error) {

	result, total, err = n.adaptors.Area.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
	}
	return
}

// Delete ...
func (n *AreaEndpoint) Delete(ctx context.Context, id int64) (err error) {

	var area *m.Area
	area, err = n.adaptors.Area.GetById(id)
	if err != nil {
		return
	}

	err = n.adaptors.Area.DeleteByName(area.Name)

	return
}

// Search ...
func (n *AreaEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Area, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Area.Search(query, limit, offset)

	return
}