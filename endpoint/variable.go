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
)

// VariableEndpoint ...
type VariableEndpoint struct {
	*CommonEndpoint
}

// NewVariableEndpoint ...
func NewVariableEndpoint(common *CommonEndpoint) *VariableEndpoint {
	return &VariableEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (v *VariableEndpoint) Add(ctx context.Context, variable m.Variable) (errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = v.validation.Valid(variable); !ok {
		return
	}

	err = v.adaptors.Variable.CreateOrUpdate(ctx, variable)

	return
}

// GetById ...
func (v *VariableEndpoint) GetById(ctx context.Context, name string) (variable m.Variable, err error) {

	variable, err = v.adaptors.Variable.GetByName(ctx, name)

	return
}

// Update ...
func (v *VariableEndpoint) Update(ctx context.Context, variable m.Variable) (errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = v.validation.Valid(variable); !ok {
		return
	}

	err = v.adaptors.Variable.CreateOrUpdate(ctx, variable)

	return
}

// GetList ...
func (v *VariableEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []m.Variable, total int64, err error) {

	list, total, err = v.adaptors.Variable.List(ctx, pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy, false)
	return
}

// Delete ...
func (v *VariableEndpoint) Delete(ctx context.Context, name string) (err error) {

	err = v.adaptors.Variable.Delete(ctx, name)
	return
}
