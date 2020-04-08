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

// MapElementEndpoint ...
type MapElementEndpoint struct {
	*CommonEndpoint
}

// NewMapElementEndpoint ...
func NewMapElementEndpoint(common *CommonEndpoint) *MapElementEndpoint {
	return &MapElementEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *MapElementEndpoint) Add(params *m.MapElement) (result *m.MapElement, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.MapElement.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.MapElement.GetById(id)

	return
}

// GetById ...
func (n *MapElementEndpoint) GetById(mId int64) (result *m.MapElement, err error) {

	result, err = n.adaptors.MapElement.GetById(mId)

	return
}

// Update ...
func (n *MapElementEndpoint) Update(params *m.MapElement) (result *m.MapElement, errs []*validation.Error, err error) {

	var m *m.MapElement
	if m, err = n.adaptors.MapElement.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&m, &params, common.JsonEngine)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.MapElement.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.MapElement.GetById(m.Id)

	_ = n.adaptors.MapZone.Clean()

	return
}

// UpdateElement ...
func (n *MapElementEndpoint) UpdateElement(params *m.MapElement) (result *m.MapElement, errs []*validation.Error, err error) {

	var old *m.MapElement
	if old, err = n.adaptors.MapElement.GetById(params.Id); err != nil {
		return
	}

	var m *m.MapElement
	common.Copy(&m, &params, common.JsonEngine)
	m.PrototypeId = old.PrototypeId
	m.PrototypeType = old.PrototypeType
	m.Prototype = old.Prototype

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.MapElement.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.MapElement.GetById(m.Id)

	_ = n.adaptors.MapZone.Clean()

	return
}

// Sort ...
func (n *MapElementEndpoint) Sort(params []*m.SortMapElement) (err error) {

	for _, s := range params {
		n.adaptors.MapElement.Sort(&m.MapElement{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

// Delete ...
func (n *MapElementEndpoint) Delete(mId int64) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	if _, err = n.adaptors.MapElement.GetById(mId); err != nil {
		return
	}

	err = n.adaptors.MapElement.Delete(mId)

	return
}

// GetList ...
func (n *MapElementEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.MapElement, total int64, err error) {

	result, total, err = n.adaptors.MapElement.List(limit, offset, order, sortBy)

	return
}
