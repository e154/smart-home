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
	"errors"
	"fmt"
	"time"

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"gorm.io/gorm/clause"

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
		DoUpdates: clause.AssignmentColumns([]string{"value", "system", "entity_id", "updated_at"}),
	}).Create(&v).Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableCreateOrUpdate)
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
			err = fmt.Errorf("%s: %w", fmt.Sprintf("name \"%s\"", name), apperr.ErrVariableNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableGet)
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
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableList)
	}
	return
}

// Delete ...
func (n Variables) Delete(ctx context.Context, name string) (err error) {
	if err = n.DB(ctx).Delete(&Variable{Name: name}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableDelete)
	}
	return
}

// List ...
func (n *Variables) List(ctx context.Context, options *adaptors.ListVariableOptions) (list []Variable, total int64, err error) {

	q := n.DB(ctx).Model(&Variable{})

	if options.System != nil {
		q = q.Where("system = ?", *options.System)
	}

	if len(options.Names) > 0 {
		q = q.Where("variables.name IN (?)", options.Names)
	}

	if options.Query != nil {
		q = q.Where("variables.name ILIKE ? OR variables.value ILIKE ?", "%"+*options.Query+"%", "%"+*options.Query+"%")
	}

	if options.Tags != nil {
		q = q.Joins(`left join variable_tags on variable_tags.variable_name = variables.name`)
		q = q.Joins(`left join tags on variable_tags.tag_id = tags.id`)
		q = q.Where("tags.name IN (?)", *options.Tags)
	}

	if options.EntityIds != nil {
		q = q.Where("variables.entity_id IN (?)", *options.EntityIds)
	}

	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableList)
		return
	}

	if options.Sort != "" && options.OrderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", options.Sort, options.OrderBy))
	}

	list = make([]Variable, 0)
	err = q.
		Preload("Tags").
		Limit(options.Limit).
		Offset(options.Offset).
		Find(&list).
		Error
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableList)
	}
	return
}

// Search ...
func (s *Variables) Search(ctx context.Context, query string, limit, offset int) (list []Variable, total int64, err error) {

	q := s.DB(ctx).Model(&Variable{}).
		Where("name ILIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableGet)
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]Variable, 0)
	if err = q.Find(&list).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableGet)
	}
	return
}

// DeleteTags ...
func (n Variables) DeleteTags(ctx context.Context, name string) (err error) {

	if err = n.DB(ctx).Model(&Variable{Name: name}).Association("Tags").Clear(); err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrVariableDeleteTag)
	}
	return
}
