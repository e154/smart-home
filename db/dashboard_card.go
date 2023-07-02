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
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common/apperr"
)

// DashboardCards ...
type DashboardCards struct {
	Db *gorm.DB
}

// DashboardCard ...
type DashboardCard struct {
	Id             int64 `gorm:"primary_key"`
	Title          string
	Weight         int
	Width          int
	Height         int
	Background     *string
	Enabled        bool
	DashboardTabId int64
	DashboardTab   *DashboardTab
	Items          []*DashboardCardItem
	Payload        json.RawMessage `gorm:"type:jsonb;not null"`
	EntityId       *common.EntityId
	Hidden         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TableName ...
func (d *DashboardCard) TableName() string {
	return "dashboard_cards"
}

// Add ...
func (n DashboardCards) Add(card *DashboardCard) (id int64, err error) {
	if err = n.Db.Create(&card).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardAdd, err.Error())
		return
	}
	id = card.Id
	return
}

// GetById ...
func (n DashboardCards) GetById(id int64) (card *DashboardCard, err error) {
	card = &DashboardCard{Id: id}
	err = n.Db.Model(card).
		Preload("Items").
		First(&card).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrDashboardCardNotFound, err.Error())
			return
		}
		err = errors.Wrap(apperr.ErrDashboardCardGet, err.Error())
	}
	return
}

// Update ...
func (n DashboardCards) Update(m *DashboardCard) (err error) {
	q := map[string]interface{}{
		"title":            m.Title,
		"height":           m.Height,
		"background":       m.Background,
		"weight":           m.Weight,
		"width":            m.Width,
		"enabled":          m.Enabled,
		"dashboard_tab_id": m.DashboardTabId,
		"payload":          m.Payload,
		"hidden":           m.Hidden,
		"entity_id":        m.EntityId,
	}

	if err = n.Db.Model(&DashboardCard{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardUpdate, err.Error())
	}
	return
}

// Delete ...
func (n DashboardCards) Delete(id int64) (err error) {
	if err = n.Db.Delete(&DashboardCard{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardDelete, err.Error())
	}
	return
}

// List ...
func (n *DashboardCards) List(limit, offset int64, orderBy, sort string) (list []*DashboardCard, total int64, err error) {

	if err = n.Db.Model(DashboardCard{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardList, err.Error())
		return
	}

	list = make([]*DashboardCard, 0)
	q := n.Db.Model(&DashboardCard{}).
		Preload("Items").
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardList, err.Error())
	}

	return
}
