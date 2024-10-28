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
	"fmt"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	pkgCommon "github.com/e154/smart-home/pkg/common"
	"github.com/pkg/errors"
)

// TelegramChats ...
type TelegramChats struct {
	*Common
}

// TelegramChat ...
type TelegramChat struct {
	EntityId  pkgCommon.EntityId
	ChatId    int64
	Username  string
	CreatedAt time.Time `gorm:"<-:create"`
}

// TableName ...
func (d *TelegramChat) TableName() string {
	return "telegram_chats"
}

// Add ...
func (n TelegramChats) Add(ctx context.Context, ch TelegramChat) (err error) {
	if err = n.DB(ctx).Create(&ch).Error; err != nil {
		err = errors.Wrap(apperr.ErrChatAdd, err.Error())
	}
	return
}

// Delete ...
func (n TelegramChats) Delete(ctx context.Context, entityId pkgCommon.EntityId, chatId int64) (err error) {
	err = n.DB(ctx).Delete(&TelegramChat{}, "entity_id = ? and chat_id = ?", entityId, chatId).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrChatDelete, err.Error())
	}
	return
}

// List ...
func (n *TelegramChats) List(ctx context.Context, limit, offset int, orderBy, sort string, entityId pkgCommon.EntityId) (list []TelegramChat, total int64, err error) {

	q := n.DB(ctx).Model(&TelegramChat{}).
		Where("entity_id = ?", entityId)

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrChatList, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	list = make([]TelegramChat, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrChatList, err.Error())
	}
	return
}
