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

	"github.com/e154/smart-home/common"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Dashboards ...
type Dashboards struct {
	Db *gorm.DB
}

// Dashboard ...
type Dashboard struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	Enabled     bool
	AreaId      *int64
	Area        *Area
	Tabs        []*DashboardTab
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Dashboard) TableName() string {
	return "dashboards"
}

// Add ...
func (n Dashboards) Add(board *Dashboard) (id int64, err error) {
	if err = n.Db.Create(&board).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = board.Id
	return
}

// GetById ...
func (n Dashboards) GetById(id int64) (board *Dashboard, err error) {
	board = &Dashboard{}
	err = n.Db.Model(board).
		Where("id = ?", id).
		Preload("Area").
		Preload("Tabs").
		Preload("Tabs.Cards").
		Preload("Tabs.Cards.Items").
		First(&board).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("id \"%s\"", id))
			return
		}
		err = errors.Wrap(err, "getById failed")
		return
	}
	return
}

// Update ...
func (n Dashboards) Update(m *Dashboard) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"enabled":     m.Enabled,
		"area_id":     m.AreaId,
	}

	if err = n.Db.Model(&Dashboard{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// Delete ...
func (n Dashboards) Delete(id int64) (err error) {
	if id == 0 {
		return
	}
	if err = n.Db.Delete(&Dashboard{Id: id}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
	}
	return
}

// List ...
func (n *Dashboards) List(limit, offset int64, orderBy, sort string) (list []*Dashboard, total int64, err error) {

	if err = n.Db.Model(Dashboard{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*Dashboard, 0)
	q := n.Db.
		Preload("Area").
		Preload("Tabs").
		Preload("Tabs.Cards").
		Preload("Tabs.Cards.Items").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(err, "find failed")
	}

	return
}