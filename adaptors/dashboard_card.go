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

package adaptors

import (
	"context"
	"strings"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

type IDashboardCard interface {
	Add(ctx context.Context, ver *m.DashboardCard) (id int64, err error)
	GetById(ctx context.Context, mapId int64) (ver *m.DashboardCard, err error)
	Update(ctx context.Context, ver *m.DashboardCard) (err error)
	Delete(ctx context.Context, id int64) (err error)
	Import(ctx context.Context, card *m.DashboardCard) (cardId int64, err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.DashboardCard, total int64, err error)
	fromDb(dbVer *db.DashboardCard) (ver *m.DashboardCard)
	toDb(ver *m.DashboardCard) (dbVer *db.DashboardCard)
}

// DashboardCard ...
type DashboardCard struct {
	IDashboardCard
	table *db.DashboardCards
	db    *gorm.DB
}

// GetDashboardCardAdaptor ...
func GetDashboardCardAdaptor(d *gorm.DB) IDashboardCard {
	return &DashboardCard{
		table: &db.DashboardCards{Db: d},
		db:    d,
	}
}

// Add ...
func (n *DashboardCard) Add(ctx context.Context, ver *m.DashboardCard) (id int64, err error) {
	dbVer := n.toDb(ver)
	id, err = n.table.Add(ctx, dbVer)
	return
}

// GetById ...
func (n *DashboardCard) GetById(ctx context.Context, mapId int64) (ver *m.DashboardCard, err error) {

	var dbVer *db.DashboardCard
	if dbVer, err = n.table.GetById(ctx, mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *DashboardCard) Update(ctx context.Context, ver *m.DashboardCard) (err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	table := db.DashboardCards{Db: tx}

	dbVer := n.toDb(ver)
	if err = table.Update(ctx, dbVer); err != nil {
		return
	}

	// items
	itemAdaptor := GetDashboardCardItemAdaptor(tx)
	for _, item := range ver.Items {
		if err = itemAdaptor.Update(ctx, item); err != nil {
			return
		}
	}

	return
}

// Delete ...
func (n *DashboardCard) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
	return
}

// Import ...
func (n *DashboardCard) Import(ctx context.Context, card *m.DashboardCard) (cardId int64, err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}
	defer func() {
		if err != nil && transaction {
			tx.Rollback()
			return
		}
		if transaction {
			err = tx.Commit().Error
		}
	}()

	cardAdaptor := GetDashboardCardAdaptor(tx)
	cardItemAdaptor := GetDashboardCardItemAdaptor(tx)

	card.Id = 0

	if !strings.Contains(card.Title, "[IMPORTED]") {
		card.Title = card.Title + " [IMPORTED]"
	}

	if cardId, err = cardAdaptor.Add(ctx, card); err != nil {
		return
	}

	// items
	if len(card.Items) > 0 {
		for _, item := range card.Items {
			item.Id = 0
			item.DashboardCardId = cardId
			if _, err = cardItemAdaptor.Add(ctx, item); err != nil {
				return
			}
		}
	}

	return
}

// List ...
func (n *DashboardCard) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.DashboardCard, total int64, err error) {

	if sort == "" {
		sort = "name"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.DashboardCard
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.DashboardCard, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

func (n *DashboardCard) fromDb(dbVer *db.DashboardCard) (ver *m.DashboardCard) {
	ver = &m.DashboardCard{
		Id:             dbVer.Id,
		Title:          dbVer.Title,
		Height:         dbVer.Height,
		Width:          dbVer.Width,
		Background:     dbVer.Background,
		Weight:         dbVer.Weight,
		Enabled:        dbVer.Enabled,
		DashboardTabId: dbVer.DashboardTabId,
		Payload:        dbVer.Payload,
		EntityId:       dbVer.EntityId,
		Hidden:         dbVer.Hidden,
		CreatedAt:      dbVer.CreatedAt,
		UpdatedAt:      dbVer.UpdatedAt,
	}

	// items
	itemAdaptor := GetDashboardCardItemAdaptor(n.db)
	for _, dbAction := range dbVer.Items {
		item := itemAdaptor.fromDb(dbAction)
		ver.Items = append(ver.Items, item)
	}

	return
}

func (n *DashboardCard) toDb(ver *m.DashboardCard) (dbVer *db.DashboardCard) {
	dbVer = &db.DashboardCard{
		Id:             ver.Id,
		Title:          ver.Title,
		Weight:         ver.Weight,
		Width:          ver.Width,
		Height:         ver.Height,
		Background:     ver.Background,
		Enabled:        ver.Enabled,
		DashboardTabId: ver.DashboardTabId,
		Payload:        ver.Payload,
		Hidden:         ver.Hidden,
		EntityId:       ver.EntityId,
	}

	return
}
