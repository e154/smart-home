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
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/e154/smart-home/common"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TelegramChats ...
type TelegramChats struct {
	Db *gorm.DB
}

// TelegramChat ...
type TelegramChat struct {
	EntityId  common.EntityId
	ChatId    int64
	Username  string
	CreatedAt time.Time
}

// TableName ...
func (d *TelegramChat) TableName() string {
	return "telegram_chats"
}

// Add ...
func (n TelegramChats) Add(ctx context.Context, ch TelegramChat) (err error) {
	if err = n.Db.WithContext(ctx).Create(&ch).Error; err != nil {
		err = errors.Wrap(apperr.ErrChatAdd, err.Error())
	}
	return
}

// Delete ...
func (n TelegramChats) Delete(ctx context.Context, entityId common.EntityId, chatId int64) (err error) {
	err = n.Db.WithContext(ctx).Delete(&TelegramChat{
		EntityId: entityId,
		ChatId:   chatId,
	}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrChatDelete, err.Error())
	}
	return
}

// List ...
func (n *TelegramChats) List(ctx context.Context, limit, offset int, orderBy, sort string, entityId common.EntityId) (list []TelegramChat, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(TelegramChat{EntityId: entityId}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrChatList, err.Error())
		return
	}

	list = make([]TelegramChat, 0)
	q := n.Db.WithContext(ctx).Model(&TelegramChat{EntityId: entityId})

	q = q.
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrChatList, err.Error())
	}
	return
}
