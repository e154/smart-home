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

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// DashboardTabs ...
type DashboardTabs struct {
	Db *gorm.DB
}

// DashboardTab ...
type DashboardTab struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Icon        string
	Enabled     bool
	Weight      int
	ColumnWidth int
	Gap         bool
	Background  *string
	DashboardId int64
	Dashboard   *Dashboard
	Cards       []*DashboardCard
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *DashboardTab) TableName() string {
	return "dashboard_tabs"
}

// Add ...
func (n DashboardTabs) Add(tab *DashboardTab) (id int64, err error) {
	if err = n.Db.Create(&tab).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabAdd, err.Error())
		return
	}
	id = tab.Id
	return
}

// GetById ...
func (n DashboardTabs) GetById(id int64) (tab *DashboardTab, err error) {
	tab = &DashboardTab{}
	err = n.Db.Model(tab).
		Where("id = ?", id).
		Preload("Cards").
		Preload("Cards.Items").
		First(&tab).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrDashboardTabNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrDashboardTabGet, err.Error())
		return
	}
	return
}

// Update ...
func (n DashboardTabs) Update(m *DashboardTab) (err error) {
	q := map[string]interface{}{
		"name":         m.Name,
		"icon":         m.Icon,
		"column_width": m.ColumnWidth,
		"gap":          m.Gap,
		"background":   m.Background,
		"enabled":      m.Enabled,
		"weight":       m.Weight,
		"dashboard_id": m.DashboardId,
	}

	if err = n.Db.Model(&DashboardTab{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabUpdate, err.Error())
	}
	return
}

// Delete ...
func (n DashboardTabs) Delete(id int64) (err error) {
	if err = n.Db.Delete(&DashboardTab{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabDelete, err.Error())
	}
	return
}

// List ...
func (n *DashboardTabs) List(limit, offset int, orderBy, sort string) (list []*DashboardTab, total int64, err error) {

	if err = n.Db.Model(DashboardTab{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabList, err.Error())
		return
	}

	list = make([]*DashboardTab, 0)
	q := n.Db.
		Preload("Cards").
		Preload("Cards.Items").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabList, err.Error())
	}

	return
}
