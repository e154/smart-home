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

package adaptors

import (
	"context"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// ITelegramChat ...
type ITelegramChat interface {
	Add(ctx context.Context, plugin m.TelegramChat) (err error)
	Delete(ctx context.Context, entityId common.EntityId, channelId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, entityId common.EntityId) (list []m.TelegramChat, total int64, err error)
	fromDb(dbVer db.TelegramChat) (ver m.TelegramChat)
	toDb(ver m.TelegramChat) (dbVer db.TelegramChat)
}

// TelegramChat ...
type TelegramChat struct {
	ITelegramChat
	table *db.TelegramChats
	db    *gorm.DB
}

// GetTelegramChannelAdaptor ...
func GetTelegramChannelAdaptor(d *gorm.DB) ITelegramChat {
	return &TelegramChat{
		table: &db.TelegramChats{Db: d},
		db:    d,
	}
}

// Add ...
func (p *TelegramChat) Add(ctx context.Context, plugin m.TelegramChat) (err error) {
	err = p.table.Add(ctx, p.toDb(plugin))
	return
}

// Delete ...
func (p *TelegramChat) Delete(ctx context.Context, entityId common.EntityId, channelId int64) (err error) {
	err = p.table.Delete(ctx, entityId, channelId)
	return
}

// List ...
func (p *TelegramChat) List(ctx context.Context, limit, offset int64, orderBy, sort string, entityId common.EntityId) (list []m.TelegramChat, total int64, err error) {
	var dbList []db.TelegramChat
	if dbList, total, err = p.table.List(ctx, int(limit), int(offset), orderBy, sort, entityId); err != nil {
		return
	}

	list = make([]m.TelegramChat, len(dbList))
	for i, dbVer := range dbList {
		list[i] = p.fromDb(dbVer)
	}
	return
}

func (p *TelegramChat) fromDb(dbVer db.TelegramChat) (ver m.TelegramChat) {
	ver = m.TelegramChat{
		EntityId:  dbVer.EntityId,
		ChatId:    dbVer.ChatId,
		Username:  dbVer.Username,
		CreatedAt: dbVer.CreatedAt,
	}

	return
}

func (p *TelegramChat) toDb(ver m.TelegramChat) (dbVer db.TelegramChat) {
	dbVer = db.TelegramChat{
		EntityId:  ver.EntityId,
		ChatId:    ver.ChatId,
		Username:  ver.Username,
		CreatedAt: ver.CreatedAt,
	}

	return
}
