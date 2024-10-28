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

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"gorm.io/gorm/clause"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Variables ...
type Variables struct {
	*Common
}

// Variable ...
type Variable struct {
	Name      string `gorm:"primary_key"`
	Value     string
	System    bool
	EntityId  *pkgCommon.EntityId
	Tags      []*Tag    `gorm:"many2many:variable_tags;"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// TableName ...
func (d *Variable) TableName() string {
	return "variables"
}

// CreateOrUpdate ...
func (n *Variables) CreateOrUpdate(ctx context.Context, v Variable) (err error) {

	err = n.DB(ctx).Omit("Tags.*").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "system", "entity_id"}),
	}).Create(&v).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrVariableCreateOrUpdate, err.Error())
	}

	return
}

// GetByName ...
func (n Variables) GetByName(ctx context.Context, name string) (variable Variable, err error) {

	variable = Variable{}
	err = n.DB(ctx).Model(&Variable{}).
		Where("name = ?", name).
		Preload("Tags").
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
	err = n.DB(ctx).Where("system = ?", true).
		Preload("Tags").
		Find(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrVariableList, err.Error())
	}
	return
}

// Delete ...
func (n Variables) Delete(ctx context.Context, name string) (err error) {
	if err = n.DB(ctx).Delete(&Variable{Name: name}).Error; err != nil {
		err = errors.Wrap(apperr.ErrVariableDelete, err.Error())
	}
	return
}

// List ...
func (n *Variables) List(ctx context.Context, limit, offset int, orderBy, sort string, system bool, name string) (list []Variable, total int64, err error) {

	q := n.DB(ctx).Model(&Variable{}).
		Preload("Tags").
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

	q := s.DB(ctx).Model(&Variable{}).
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

// DeleteTags ...
func (n Variables) DeleteTags(ctx context.Context, name string) (err error) {

	if err = n.DB(ctx).Model(&Variable{Name: name}).Association("Tags").Clear(); err != nil {
		err = errors.Wrap(apperr.ErrVariableDeleteTag, err.Error())
	}
	return
}
