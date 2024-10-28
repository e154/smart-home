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

	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/models"
)

// TemplateEndpoint ...
type TemplateEndpoint struct {
	*CommonEndpoint
}

// NewTemplateEndpoint ...
func NewTemplateEndpoint(common *CommonEndpoint) *TemplateEndpoint {
	return &TemplateEndpoint{
		CommonEndpoint: common,
	}
}

// UpdateOrCreate ...
func (t *TemplateEndpoint) UpdateOrCreate(ctx context.Context, params *models.Template) (err error) {

	if ok, errs := t.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	err = t.adaptors.Template.UpdateOrCreate(ctx, params)
	return
}

// UpdateStatus ...
func (t *TemplateEndpoint) UpdateStatus(ctx context.Context, params *models.Template) (err error) {

	if ok, errs := t.validation.Valid(params); !ok {
		err = apperr.ErrValidation
		apperr.SetValidationErrors(err, errs)
		return
	}

	err = t.adaptors.Template.UpdateStatus(ctx, params)
	return
}

// GetByName ...
func (t *TemplateEndpoint) GetByName(ctx context.Context, name string) (result *models.Template, err error) {
	result, err = t.adaptors.Template.GetByName(ctx, name)
	if err != nil {
		return
	}

	err = t.adaptors.Template.GetMarkers(ctx, result)
	return
}

// GetItemByName ...
func (t *TemplateEndpoint) GetItemByName(ctx context.Context, name string) (result *models.Template, err error) {
	result, err = t.adaptors.Template.GetItemByName(ctx, name)
	if err != nil {
		return
	}
	return
}

// GetItemsSortedList ...
func (t *TemplateEndpoint) GetItemsSortedList(ctx context.Context) (count int64, items []string, err error) {
	count, items, err = t.adaptors.Template.GetItemsSortedList(ctx)

	return
}

// GetList ...
func (t *TemplateEndpoint) GetList(ctx context.Context) (count int64, templates []*models.Template, err error) {
	templates, err = t.adaptors.Template.GetList(ctx, models.TemplateTypeTemplate)

	return
}

// Delete ...
func (t *TemplateEndpoint) Delete(ctx context.Context, name string) (err error) {
	err = t.adaptors.Template.Delete(ctx, name)
	return
}

// GetItemsTree ...
func (t *TemplateEndpoint) GetItemsTree(ctx context.Context) (tree []*models.TemplateTree, err error) {
	tree, err = t.adaptors.Template.GetItemsTree(ctx)
	return
}

// UpdateItemsTree ...
func (t *TemplateEndpoint) UpdateItemsTree(ctx context.Context, tree []*models.TemplateTree) (err error) {
	err = t.adaptors.Template.UpdateItemsTree(ctx, tree)
	return
}

// Search ...
func (t *TemplateEndpoint) Search(ctx context.Context, query string, limit, offset int) (result []*models.Template, total int64, err error) {
	result, total, err = t.adaptors.Template.Search(ctx, query, limit, offset)
	return
}

// Preview ...
func (t *TemplateEndpoint) Preview(ctx context.Context, template *models.TemplateContent) (data string, err error) {

	var items []*models.Template
	if items, err = t.adaptors.Template.GetList(ctx, models.TemplateTypeItem); err != nil {
		return
	}

	data, err = models.PreviewTemplate(items, template)

	return
}
