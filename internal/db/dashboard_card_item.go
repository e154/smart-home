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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// DashboardCardItems ...
type DashboardCardItems struct {
	*Common
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
	EntityId        *pkgCommon.EntityId
	Payload         json.RawMessage `gorm:"type:jsonb;not null"`
	Hidden          bool
	Frozen          bool
	CreatedAt       time.Time `gorm:"<-:create"`
	UpdatedAt       time.Time
}

// TableName ...
func (d *DashboardCardItem) TableName() string {
	return "dashboard_card_items"
}

// Add ...
func (n DashboardCardItems) Add(ctx context.Context, item *DashboardCardItem) (id int64, err error) {
	if err = n.DB(ctx).Create(&item).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation:
				if strings.Contains(pgErr.Message, "dashboard_card_item_2_entities_fk") {
					details := pgErr.Detail
					details = strings.Split(details, `Key (entity_id)=(`)[1]
					details = strings.Split(details, `) is not present in table "entities".`)[0]
					err = fmt.Errorf("%s: %w", fmt.Sprintf("with name \"%s\"", details), apperr.ErrEntityGet)
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrDashboardCardItemAdd)
		return
	}
	id = item.Id
	return
}

// GetById ...
func (n DashboardCardItems) GetById(ctx context.Context, id int64) (item *DashboardCardItem, err error) {
	item = &DashboardCardItem{Id: id}
	if err = n.DB(ctx).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("%s: %w", fmt.Sprintf("with id \"%d\"", id), apperr.ErrDashboardCardItemNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrDashboardCardItemGet)
	}
	return
}

// Update ...
func (n DashboardCardItems) Update(ctx context.Context, m *DashboardCardItem) (err error) {
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

	if err = n.DB(ctx).Model(&DashboardCardItem{Id: m.Id}).Updates(q).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrDashboardCardItemUpdate)
	}
	return
}

// Delete ...
func (n DashboardCardItems) Delete(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Delete(&DashboardCardItem{Id: id}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrDashboardCardItemDelete)
	}
	return
}

// List ...
func (n *DashboardCardItems) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*DashboardCardItem, total int64, err error) {

	if err = n.DB(ctx).Model(DashboardCardItem{}).Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrDashboardCardItemList)
		return
	}

	list = make([]*DashboardCardItem, 0)
	q := n.DB(ctx).Model(&DashboardCardItem{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrDashboardCardItemList)
	}

	return
}
