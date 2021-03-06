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
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

// EntityEndpoint ...
type EntityEndpoint struct {
	*CommonEndpoint
}

// NewEntityEndpoint ...
func NewEntityEndpoint(common *CommonEndpoint) *EntityEndpoint {
	return &EntityEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *EntityEndpoint) Add(entity *m.Entity) (result *m.Entity, errs []*validation.Error, err error) {

	_, errs = entity.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Entity.Add(entity); err != nil {
		return
	}

	if result, err = n.adaptors.Entity.GetById(entity.Id); err != nil {
		return
	}

	err = n.entityManager.Add(entity)

	return
}

// GetById ...
func (n *EntityEndpoint) GetById(id common.EntityId) (result *m.Entity, err error) {

	result, err = n.adaptors.Entity.GetById(id)

	return
}

// Update ...
func (n *EntityEndpoint) Update(params *m.Entity) (result *m.Entity, errs []*validation.Error, err error) {

	var entity *m.Entity
	if entity, err = n.adaptors.Entity.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&entity, &params, common.JsonEngine)

	// validation
	_, errs = entity.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Entity.Update(entity); err != nil {
		return
	}

	if result, err = n.adaptors.Entity.GetById(entity.Id); err != nil {
		return
	}

	err = n.entityManager.Update(entity)

	return
}

// List ...
func (n *EntityEndpoint) List(limit, offset int64, order, sortBy string) (result []*m.Entity, total int64, err error) {
	result, total, err = n.adaptors.Entity.List(limit, offset, order, sortBy, false)
	return
}

// Delete ...
func (n *EntityEndpoint) Delete(id common.EntityId) (err error) {

	if id == "" {
		err = errors.New("entity id is null")
		return
	}

	var entity *m.Entity
	if entity, err = n.adaptors.Entity.GetById(id); err != nil {
		return
	}

	if err = n.adaptors.Entity.Delete(entity.Id); err != nil {
		return
	}

	n.entityManager.Remove(id)

	return
}

// Search ...
func (n *EntityEndpoint) Search(query string, limit, offset int) (result []*m.Entity, total int64, err error) {

	result, total, err = n.adaptors.Entity.Search(query, limit, offset)

	return
}
