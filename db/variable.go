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

package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
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
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// TableName ...
func (d *Variable) TableName() string {
	return "variables"
}

// Add ...
func (n Variables) Add(ctx context.Context, variable Variable) (err error) {
	if err = n.Db.WithContext(ctx).Create(&variable).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableAdd, err.Error())
	}
	return
}

// CreateOrUpdate ...
func (n *Variables) CreateOrUpdate(ctx context.Context, v Variable) (err error) {
	params := map[string]interface{}{
		"name":  v.Name,
		"value": v.Value,
	}
	if n.Db.WithContext(ctx).Model(&v).Where("name = ?", v.Name).Updates(params).RowsAffected == 0 {
		err = n.Db.WithContext(ctx).Create(&v).Error
	}
	return
}

// GetByName ...
func (n Variables) GetByName(ctx context.Context, name string) (variable Variable, err error) {
	variable = Variable{}
	err = n.Db.WithContext(ctx).Model(&Variable{}).
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
func (n Variables) GetAllSystem(ctx context.Context) (list []Variable, err error) {
	list = make([]Variable, 0)
	err = n.Db.WithContext(ctx).Where("system = ?", true).
		Find(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
	}
	return
}

// Update ...
func (n Variables) Update(ctx context.Context, m Variable) (err error) {
	err = n.Db.WithContext(ctx).Model(&Variable{Name: m.Name}).Updates(map[string]interface{}{
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
func (n Variables) Delete(ctx context.Context, name string) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&Variable{Name: name}).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableDelete, err.Error())
	}
	return
}

// List ...
func (n *Variables) List(ctx context.Context, limit, offset int, orderBy, sort string, system bool, name string) (list []Variable, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Variable{}).
		Where("system = ?", system)

	if strings.Contains(name, ",") {
		names := strings.Split(name, ",")
		if len(names) > 0 {
			q = q.Where("name IN (?)", names)
		}
	}

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
		return
	}

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	list = make([]Variable, 0)
	err = q.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
	}
	return
}

// Search ...
func (s *Variables) Search(ctx context.Context, query string, limit, offset int) (list []Variable, total int64, err error) {

	q := s.Db.WithContext(ctx).Model(&Variable{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableGet, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]Variable, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableGet, err.Error())
	}
	return
}
