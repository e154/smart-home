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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/e154/smart-home/common/apperr"
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
func (n Dashboards) Add(ctx context.Context, board *Dashboard) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&board).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_dashboards_unq") {
					err = errors.Wrap(apperr.ErrDashboardAdd, fmt.Sprintf("dashboard name \"%s\" not unique", board.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrDashboardAdd, err.Error())
		return
	}
	id = board.Id
	return
}

// GetById ...
func (n Dashboards) GetById(ctx context.Context, id int64) (board *Dashboard, err error) {
	board = &Dashboard{}
	err = n.Db.WithContext(ctx).Model(board).
		Where("id = ?", id).
		Preload("Area").
		Preload("Tabs").
		Preload("Tabs.Cards").
		Preload("Tabs.Cards.Items").
		First(&board).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrDashboardNotFound, fmt.Sprintf("id \"%s\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrDashboardGet, err.Error())
		return
	}
	return
}

// Update ...
func (n Dashboards) Update(ctx context.Context, board *Dashboard) (err error) {
	q := map[string]interface{}{
		"name":        board.Name,
		"description": board.Description,
		"enabled":     board.Enabled,
		"area_id":     board.AreaId,
	}

	if err = n.Db.WithContext(ctx).Model(&Dashboard{Id: board.Id}).Updates(q).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_dashboards_unq") {
					err = errors.Wrap(apperr.ErrDashboardUpdate, fmt.Sprintf("dashboard name \"%s\" not unique", board.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrDashboardUpdate, err.Error())
	}
	return
}

// Delete ...
func (n Dashboards) Delete(ctx context.Context, id int64) (err error) {
	if id == 0 {
		return
	}
	if err = n.Db.WithContext(ctx).Delete(&Dashboard{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardDelete, err.Error())
	}
	return
}

// List ...
func (n *Dashboards) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*Dashboard, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(Dashboard{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardGet, err.Error())
		return
	}

	list = make([]*Dashboard, 0)
	q := n.Db.WithContext(ctx).
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
		err = errors.Wrap(apperr.ErrDashboardList, err.Error())
	}

	return
}

// Search ...
func (d *Dashboards) Search(ctx context.Context, query string, limit, offset int) (list []*Dashboard, total int64, err error) {

	q := d.Db.Model(&Dashboard{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Dashboard, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrDashboardSearch, err.Error())
	}
	return
}
