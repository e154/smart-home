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
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// DashboardTabs ...
type DashboardTabs struct {
	*Common
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
	Payload     json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time       `gorm:"<-:create"`
	UpdatedAt   time.Time
}

// TableName ...
func (d *DashboardTab) TableName() string {
	return "dashboard_tabs"
}

// Add ...
func (n DashboardTabs) Add(ctx context.Context, tab *DashboardTab) (id int64, err error) {
	if err = n.DB(ctx).Create(&tab).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_dashboard_tabs_unq") {
					err = errors.Wrap(apperr.ErrDashboardTabAdd, fmt.Sprintf("tab name \"%s\" not unique", tab.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrDashboardTabAdd, err.Error())
		return
	}
	id = tab.Id
	return
}

// GetById ...
func (n DashboardTabs) GetById(ctx context.Context, id int64) (tab *DashboardTab, err error) {
	tab = &DashboardTab{}
	err = n.DB(ctx).Model(tab).
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
func (n DashboardTabs) Update(ctx context.Context, tab *DashboardTab) (err error) {
	q := map[string]interface{}{
		"name":         tab.Name,
		"icon":         tab.Icon,
		"column_width": tab.ColumnWidth,
		"gap":          tab.Gap,
		"background":   tab.Background,
		"enabled":      tab.Enabled,
		"weight":       tab.Weight,
		"dashboard_id": tab.DashboardId,
		"payload":      tab.Payload,
	}

	if err = n.DB(ctx).Model(&DashboardTab{Id: tab.Id}).Updates(q).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_dashboard_tabs_unq") {
					err = errors.Wrap(apperr.ErrDashboardTabUpdate, fmt.Sprintf("tab name \"%s\" not unique", tab.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrDashboardTabUpdate, err.Error())
	}
	return
}

// Delete ...
func (n DashboardTabs) Delete(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Delete(&DashboardTab{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabDelete, err.Error())
	}
	return
}

// List ...
func (n *DashboardTabs) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*DashboardTab, total int64, err error) {

	if err = n.DB(ctx).Model(DashboardTab{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardTabList, err.Error())
		return
	}

	list = make([]*DashboardTab, 0)
	q := n.DB(ctx).
		WithContext(ctx).
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
