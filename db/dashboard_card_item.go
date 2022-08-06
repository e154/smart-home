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
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
)

// DashboardCardItems ...
type DashboardCardItems struct {
	Db *gorm.DB
}

// DashboardCardItem ...
type DashboardCardItem struct {
	Id              int64 `gorm:"primary_key"`
	Title           string
	Type            string
	Weight          int
	Enabled         bool
	DashboardCardId int64
	DashboardCard   *DashboardCard
	EntityId        *common.EntityId
	Payload         json.RawMessage `gorm:"type:jsonb;not null"`
	Hidden          bool
	Frozen          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// TableName ...
func (d *DashboardCardItem) TableName() string {
	return "dashboard_card_items"
}

// Add ...
func (n DashboardCardItems) Add(item *DashboardCardItem) (id int64, err error) {
	if err = n.Db.Create(&item).Error; err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				if strings.Contains(pgErr.Message, "dashboard_card_item_2_entities_fk") {
					details := pgErr.Detail
					details = strings.Split(details, `Key (entity_id)=(`)[1]
					details = strings.Split(details, `) is not present in table "entities".`)[0]
					err = errors.Wrap(apperr.ErrEntityGet, fmt.Sprintf("with name \"%s\"", details))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrDashboardCardItemAdd, err.Error())
		return
	}
	id = item.Id
	return
}

// GetById ...
func (n DashboardCardItems) GetById(id int64) (item *DashboardCardItem, err error) {
	item = &DashboardCardItem{Id: id}
	if err = n.Db.First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrDashboardCardItemNotFound, fmt.Sprintf("with id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrDashboardCardItemGet, err.Error())
	}
	return
}

// Update ...
func (n DashboardCardItems) Update(m *DashboardCardItem) (err error) {
	q := map[string]interface{}{
		"title":             m.Title,
		"type":              m.Type,
		"weight":            m.Weight,
		"enabled":           m.Enabled,
		"dashboard_card_id": m.DashboardCardId,
		"entity_id":         m.EntityId,
		"payload":           m.Payload,
		"hidden":            m.Hidden,
	}

	if err = n.Db.Model(&DashboardCardItem{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardItemUpdate, err.Error())
	}
	return
}

// Delete ...
func (n DashboardCardItems) Delete(id int64) (err error) {
	if err = n.Db.Delete(&DashboardCardItem{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardItemDelete, err.Error())
	}
	return
}

// List ...
func (n *DashboardCardItems) List(limit, offset int64, orderBy, sort string) (list []*DashboardCardItem, total int64, err error) {

	if err = n.Db.Model(DashboardCardItem{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardItemList, err.Error())
		return
	}

	list = make([]*DashboardCardItem, 0)
	q := n.Db.Model(&DashboardCardItem{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardCardItemList, err.Error())
	}

	return
}
