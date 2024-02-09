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
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strings"

	"github.com/e154/smart-home/common/apperr"
)

// Tags ...
type Tags struct {
	Db *gorm.DB
}

// Tag ...
type Tag struct {
	Id   int64 `gorm:"primary_key"`
	Name string
}

// TableName ...
func (d *Tag) TableName() string {
	return "tags"
}

// Add ...
func (n Tags) Add(ctx context.Context, tag *Tag) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(tag).Error; err != nil {
		err = errors.Wrap(apperr.ErrTagAdd, err.Error())
		return
	}
	id = tag.Id
	return
}

// List ...
func (n *Tags) List(ctx context.Context, limit, offset int, orderBy, sort string, query *string, names *[]string) (list []*Tag, total int64, err error) {

	list = make([]*Tag, 0)
	q := n.Db.WithContext(ctx).Model(Tag{})
	if query != nil {
		q = q.Where("name LIKE ? or source LIKE ?", "%"+*query+"%", "%"+*query+"%")
	}
	if names != nil {
		q = q.Where("name IN (?)", *names)
	}
	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTagList, err.Error())
		return
	}
	err = q.
		Limit(limit).
		Offset(offset).
		//Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrTagList, err.Error())
	}
	return
}

// GetByName ...
func (n *Tags) GetByName(ctx context.Context, name string) (tag *Tag, err error) {
	tag = &Tag{}
	err = n.Db.WithContext(ctx).Model(tag).
		Where("name = ?", name).
		First(&tag).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrTagNotFound, fmt.Sprintf("name \"%s\"", name))
			return
		}
		err = errors.Wrap(apperr.ErrTagGet, err.Error())
	}

	return
}

// GetById ...
func (n *Tags) GetById(ctx context.Context, id int64) (tag *Tag, err error) {
	tag = &Tag{}
	err = n.Db.WithContext(ctx).Model(tag).
		Where("id = ?", id).
		First(&tag).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrTagNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrTagGet, err.Error())
	}

	return
}

// Delete ...
func (n *Tags) Delete(ctx context.Context, name string) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&Tag{Name: name}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTagDelete, err.Error())
	}
	return
}

// Update ...
func (n *Tags) Update(ctx context.Context, tag *Tag) (err error) {
	err = n.Db.WithContext(ctx).Model(&Tag{Id: tag.Id}).Updates(map[string]interface{}{
		"name": tag.Name,
	}).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "tag_name_unq") {
					err = errors.Wrap(apperr.ErrTagUpdate, fmt.Sprintf("tag name \"%s\" not unique", tag.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrTagUpdate, err.Error())
		return
	}

	return
}

// Search ...
func (n *Tags) Search(ctx context.Context, query string, limit, offset int) (list []*Tag, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Tag{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTagSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Tag, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrTagSearch, err.Error())
	}
	return
}
