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

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
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
func (v *VariableEndpoint) Add(ctx context.Context, variable m.Variable) (err error) {

	if ok, errs := v.validation.Valid(variable); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	if err = v.adaptors.Variable.CreateOrUpdate(ctx, variable); err != nil {
		return
	}

	v.eventBus.Publish(fmt.Sprintf("system/models/variables/%s", variable.Name), events.EventUpdatedVariableModel{
		Name:  variable.Name,
		Value: variable.Value,
	})

	return
}

// GetById ...
func (v *VariableEndpoint) GetById(ctx context.Context, name string) (variable m.Variable, err error) {

	variable, err = v.adaptors.Variable.GetByName(ctx, name)

	return
}

// Update ...
func (v *VariableEndpoint) Update(ctx context.Context, _variable m.Variable) (err error) {

	if ok, errs := v.validation.Valid(_variable); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetValidationErrors(err, errs)
		return
	}

	var variable m.Variable
	if variable, err = v.adaptors.Variable.GetByName(ctx, _variable.Name); err == nil {
		if variable.System && v.checkSuperUser(ctx) {
			err = apperr.ErrVariableUpdateForbidden
			return
		}
		variable.Value = _variable.Value
	} else {
		variable.Name = _variable.Name
		variable.Value = _variable.Value
	}

	if err = v.adaptors.Variable.CreateOrUpdate(ctx, variable); err != nil {
		return
	}

	v.eventBus.Publish(fmt.Sprintf("system/models/variables/%s", variable.Name), events.EventUpdatedVariableModel{
		Name:  variable.Name,
		Value: variable.Value,
	})

	return
}

// GetList ...
func (v *VariableEndpoint) GetList(ctx context.Context, pagination common.PageParams) (list []m.Variable, total int64, err error) {

	list, total, err = v.adaptors.Variable.List(ctx, pagination.Limit, pagination.Offset, pagination.Order,
		pagination.SortBy, false, "")
	return
}

// Delete ...
func (v *VariableEndpoint) Delete(ctx context.Context, name string) (err error) {

	var variable m.Variable
	if variable, err = v.adaptors.Variable.GetByName(ctx, name); err == nil {
		if variable.System && v.checkSuperUser(ctx) {
			err = apperr.ErrVariableUpdateForbidden
			return
		}

	}

	if err = v.adaptors.Variable.Delete(ctx, name); err != nil {
		return
	}

	v.eventBus.Publish(fmt.Sprintf("system/models/variables/%s", name), events.EventRemovedVariableModel{
		Name: name,
	})

	return
}

// Search ...
func (n *VariableEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []m.Variable, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Variable.Search(ctx, query, int(limit), int(offset))

	return
}
