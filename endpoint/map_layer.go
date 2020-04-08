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

// MapLayerEndpoint ...
type MapLayerEndpoint struct {
	*CommonEndpoint
}

// NewMapLayerEndpoint ...
func NewMapLayerEndpoint(common *CommonEndpoint) *MapLayerEndpoint {
	return &MapLayerEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *MapLayerEndpoint) Add(params *m.MapLayer) (result *m.MapLayer, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.MapLayer.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.MapLayer.GetById(id)

	return
}

// GetById ...
func (n *MapLayerEndpoint) GetById(mId int64) (result *m.MapLayer, err error) {

	result, err = n.adaptors.MapLayer.GetById(mId)

	return
}

// Update ...
func (n *MapLayerEndpoint) Update(params *m.MapLayer) (result *m.MapLayer, errs []*validation.Error, err error) {

	var m *m.MapLayer
	if m, err = n.adaptors.MapLayer.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&m, &params, common.JsonEngine)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.MapLayer.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.MapLayer.GetById(m.Id)

	return
}

// Sort ...
func (n *MapLayerEndpoint) Sort(params []*m.SortMapLayer) (err error) {

	for _, s := range params {
		n.adaptors.MapLayer.Sort(&m.MapLayer{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

// Delete ...
func (n *MapLayerEndpoint) Delete(mId int64) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	if _, err = n.adaptors.MapLayer.GetById(mId); err != nil {
		return
	}

	err = n.adaptors.MapLayer.Delete(mId)

	return
}

// GetList ...
func (n *MapLayerEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.MapLayer, total int64, err error) {

	result, total, err = n.adaptors.MapLayer.List(limit, offset, order, sortBy)

	return
}
