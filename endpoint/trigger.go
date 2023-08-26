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

// TriggerEndpoint ...
type TriggerEndpoint struct {
	*CommonEndpoint
}

// NewTriggerEndpoint ...
func NewTriggerEndpoint(common *CommonEndpoint) *TriggerEndpoint {
	return &TriggerEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *TriggerEndpoint) Add(ctx context.Context, trigger *m.Trigger) (result *m.Trigger, errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = n.validation.Valid(trigger); !ok {
		err = apperr.ErrInvalidRequest
		apperr.SetContext(err, errs)
		return
	}

	if trigger.Id, err = n.adaptors.Trigger.Add(trigger); err != nil {
		return
	}

	result, err = n.adaptors.Trigger.GetById(trigger.Id)

	return
}

// GetById ...
func (n *TriggerEndpoint) GetById(ctx context.Context, id int64) (result *m.Trigger, err error) {

	result, err = n.adaptors.Trigger.GetById(id)

	return
}

// Update ...
func (n *TriggerEndpoint) Update(ctx context.Context, params *m.Trigger) (result *m.Trigger, errs validator.ValidationErrorsTranslations, err error) {

	_, err = n.adaptors.Trigger.GetById(params.Id)
	if err != nil {
		return
	}


	var ok bool
	if ok, errs = n.validation.Valid(params); !ok {
		return
	}

	if err = n.adaptors.Trigger.Update(params); err != nil {
		return
	}

	result, err = n.adaptors.Trigger.GetById(params.Id)

	return
}

// GetList ...
func (n *TriggerEndpoint) GetList(ctx context.Context, pagination common.PageParams) (result []*m.Trigger, total int64, err error) {

	result, total, err = n.adaptors.Trigger.List(pagination.Limit, pagination.Offset, pagination.Order, pagination.SortBy)
	if err != nil {
	}
	return
}

// Delete ...
func (n *TriggerEndpoint) Delete(ctx context.Context, id int64) (err error) {

	var trigger *m.Trigger
	trigger, err = n.adaptors.Trigger.GetById(id)
	if err != nil {
		return
	}

	err = n.adaptors.Trigger.Delete(trigger.Id)

	return
}

// Search ...
func (n *TriggerEndpoint) Search(ctx context.Context, query string, limit, offset int64) (result []*m.Trigger, total int64, err error) {

	if limit == 0 {
		limit = common.DefaultPageSize
	}

	result, total, err = n.adaptors.Trigger.Search(query, int(limit), int(offset))

	return
}
