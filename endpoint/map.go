// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

// MapEndpoint ...
type MapEndpoint struct {
	*CommonEndpoint
}

// NewMapEndpoint ...
func NewMapEndpoint(common *CommonEndpoint) *MapEndpoint {
	return &MapEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (m *MapEndpoint) Add(params *m.Map) (result *m.Map, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = m.adaptors.Map.Add(params); err != nil {
		return
	}

	result, err = m.adaptors.Map.GetById(id)

	return
}

// GetById ...
func (m *MapEndpoint) GetById(id int64) (result *m.Map, err error) {

	result, err = m.adaptors.Map.GetById(id)

	return
}

// GetActiveElements ...
func (m *MapEndpoint) GetActiveElements(sortBy, order string, limit, offset int) (result []*m.MapElement, total int64, err error) {

	result, total, err = m.adaptors.MapElement.GetActiveElements(sortBy, order, limit, offset)

	return
}

// GetFullById ...
func (m *MapEndpoint) GetFullById(mId int64) (result *m.Map, err error) {

	result, err = m.adaptors.Map.GetFullById(mId)

	return
}

// Update ...
func (n *MapEndpoint) Update(params *m.Map) (result *m.Map, errs []*validation.Error, err error) {

	var m *m.Map
	if m, err = n.adaptors.Map.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&m, &params, common.JsonEngine)

	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Map.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.Map.GetById(params.Id)

	return
}

// GetList ...
func (n *MapEndpoint) GetList(limit, offset int64, order, sortBy string) (items []*m.Map, total int64, err error) {

	items, total, err = n.adaptors.Map.List(limit, offset, order, sortBy)

	return
}

// Delete ...
func (n *MapEndpoint) Delete(mId int64) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	var m *m.Map
	if m, err = n.adaptors.Map.GetById(mId); err != nil {
		return
	}

	err = n.adaptors.Map.Delete(m.Id)

	return
}

// Search ...
func (n *MapEndpoint) Search(query string, limit, offset int) (items []*m.Map, total int64, err error) {

	items, total, err = n.adaptors.Map.Search(query, limit, offset)

	return
}
