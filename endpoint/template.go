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
	m "github.com/e154/smart-home/models"
	"github.com/go-playground/validator/v10"
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
func (t *TemplateEndpoint) UpdateOrCreate(params *m.Template) (errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = t.validation.Valid(params); !ok {
		return
	}

	err = t.adaptors.Template.UpdateOrCreate(params)
	return
}

// UpdateStatus ...
func (t *TemplateEndpoint) UpdateStatus(params *m.Template) (errs validator.ValidationErrorsTranslations, err error) {

	var ok bool
	if ok, errs = t.validation.Valid(params); !ok {
		return
	}

	err = t.adaptors.Template.UpdateStatus(params)
	return
}

// GetByName ...
func (t *TemplateEndpoint) GetByName(name string) (result *m.Template, err error) {
	result, err = t.adaptors.Template.GetByName(name)
	if err != nil {
		return
	}

	err = t.adaptors.Template.GetMarkers(result)
	return
}

// GetItemByName ...
func (t *TemplateEndpoint) GetItemByName(name string) (result *m.Template, err error) {
	result, err = t.adaptors.Template.GetItemByName(name)
	if err != nil {
		return
	}
	return
}

// GetItemsSortedList ...
func (t *TemplateEndpoint) GetItemsSortedList() (count int64, items []string, err error) {
	count, items, err = t.adaptors.Template.GetItemsSortedList()

	return
}

// GetList ...
func (t *TemplateEndpoint) GetList() (count int64, templates []*m.Template, err error) {
	templates, err = t.adaptors.Template.GetList(m.TemplateTypeTemplate)

	return
}

// Delete ...
func (t *TemplateEndpoint) Delete(name string) (err error) {
	err = t.adaptors.Template.Delete(name)
	return
}

// GetItemsTree ...
func (t *TemplateEndpoint) GetItemsTree() (tree []*m.TemplateTree, err error) {
	tree, err = t.adaptors.Template.GetItemsTree()
	return
}

// UpdateItemsTree ...
func (t *TemplateEndpoint) UpdateItemsTree(tree []*m.TemplateTree) (err error) {
	err = t.adaptors.Template.UpdateItemsTree(tree)
	return
}

// Search ...
func (t *TemplateEndpoint) Search(query string, limit, offset int) (result []*m.Template, total int64, err error) {
	result, total, err = t.adaptors.Template.Search(query, limit, offset)
	return
}

// Preview ...
func (t *TemplateEndpoint) Preview(template *m.TemplateContent) (data string, err error) {

	var items []*m.Template
	if items, err = t.adaptors.Template.GetList(m.TemplateTypeItem); err != nil {
		return
	}

	data, err = m.PreviewTemplate(items, template)

	return
}
