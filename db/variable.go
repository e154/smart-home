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

package db

import (
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Variables ...
type Variables struct {
	Db *gorm.DB
}

// Variable ...
type Variable struct {
	Name      string `gorm:"primary_key"`
	Value     string
	System    bool
	EntityId  *common.EntityId
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName ...
func (d *Variable) TableName() string {
	return "variables"
}

// Add ...
func (n Variables) Add(variable Variable) (err error) {
	if err = n.Db.Create(&variable).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableAdd, err.Error())
	}
	return
}

// CreateOrUpdate ...
func (n *Variables) CreateOrUpdate(v Variable) (err error) {
	var entityId = "null"
	if v.EntityId != nil {
		entityId = v.EntityId.String()
	}
	err = n.Db.Model(&Variable{}).
		Set("gorm:insert_option",
			fmt.Sprintf("ON CONFLICT (name) DO UPDATE SET value = '%s', entity_id = %s, updated_at = '%s'", v.Value, entityId, time.Now().Format(time.RFC3339))).
		Create(&v).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableUpdate, err.Error())
	}
	return
}

// GetByName ...
func (n Variables) GetByName(name string) (variable Variable, err error) {
	variable = Variable{}
	err = n.Db.Model(&Variable{}).
		Where("name = ?", name).
		First(&variable).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrVariableNotFound, fmt.Sprintf("name \"%s\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrVariableGet, err.Error())
	}
	return
}

// GetAllSystem ...
func (n Variables) GetAllSystem() (list []Variable, err error) {
	list = make([]Variable, 0)
	err = n.Db.Where("system = ?", true).
		Find(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
	}
	return
}

// Update ...
func (n Variables) Update(m Variable) (err error) {
	err = n.Db.Model(&Variable{Name: m.Name}).Updates(map[string]interface{}{
		"value":     m.Value,
		"system":    m.System,
		"entity_id": m.EntityId,
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableUpdate, err.Error())
	}
	return
}

// Delete ...
func (n Variables) Delete(name string) (err error) {
	if err = n.Db.Delete(&Variable{Name: name}).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableDelete, err.Error())
	}
	return
}

// List ...
func (n *Variables) List(limit, offset int64, orderBy, sort string, system bool) (list []Variable, total int64, err error) {

	if err = n.Db.Model(Variable{}).Where("system = ?", system).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
		return
	}

	list = make([]Variable, 0)
	err = n.Db.
		Model(&Variable{}).
		Where("system = ?", system).
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
	}
	return
}
