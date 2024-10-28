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
	"gorm.io/gorm"
)

// MessageDeliveries ...
type MessageDeliveries struct {
	*Common
}

// MessageDelivery ...
type MessageDelivery struct {
	Id                 int64 `gorm:"primary_key"`
	Message            *Message
	MessageId          int64
	Address            string
	EntityId           *pkgCommon.EntityId
	Status             string
	ErrorMessageStatus *string   `gorm:"column:error_system_code"`
	ErrorMessageBody   *string   `gorm:"column:error_system_message"`
	CreatedAt          time.Time `gorm:"<-:create"`
	UpdatedAt          time.Time
}

// TableName ...
func (d *MessageDelivery) TableName() string {
	return "message_deliveries"
}

// Add ...
func (n *MessageDeliveries) Add(ctx context.Context, msg *MessageDelivery) (id int64, err error) {
	if err = n.DB(ctx).Create(&msg).Error; err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryAdd, err.Error())
		return
	}
	id = msg.Id
	return
}

func (n *MessageDeliveries) List(ctx context.Context, limit, offset int, orderBy, sort string, queryObj *MessageDeliveryQuery) (list []*MessageDelivery, total int64, err error) {

	list = make([]*MessageDelivery, 0)
	q := n.DB(ctx).Model(&MessageDelivery{}).
		Joins(`left join messages on messages.id = message_deliveries.message_id`)

	if sort != "" && orderBy != "" {
		q = q.Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	q = q.Preload("Message")
	if queryObj != nil {
		if queryObj.StartDate != nil {
			q = q.Where("message_deliveries.created_at >= ?", &queryObj.StartDate)
		}
		if queryObj.EndDate != nil {
			q = q.Where("message_deliveries.created_at <= ?", &queryObj.EndDate)
		}
		if len(queryObj.Types) > 0 {
			q = q.Where("messages.type in (?)", queryObj.Types)
		}
	}

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryList, err.Error())
		return
	}

	err = q.
		Limit(limit).
		Offset(offset).
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryList, err.Error())
	}
	return
}

// GetAllUncompleted ...
func (n *MessageDeliveries) GetAllUncompleted(ctx context.Context, limit, offset int) (list []*MessageDelivery, total int64, err error) {

	if err = n.DB(ctx).Model(&MessageDelivery{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryUpdate, err.Error())
		return
	}

	list = make([]*MessageDelivery, 0)
	err = n.DB(ctx).
		Where("status in ('in_progress', 'new')").
		Limit(limit).
		Offset(offset).
		Preload("Message").
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryUpdate, err.Error())
	}
	return
}

// SetStatus ...
func (n *MessageDeliveries) SetStatus(ctx context.Context, msg *MessageDelivery) (err error) {

	err = n.DB(ctx).Model(&MessageDelivery{Id: msg.Id}).
		Updates(map[string]interface{}{
			"status":               msg.Status,
			"error_system_code":    msg.ErrorMessageStatus,
			"error_system_message": msg.ErrorMessageBody,
		}).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryUpdate, err.Error())
	}
	return
}

// Delete ...
func (n *MessageDeliveries) Delete(ctx context.Context, id int64) (err error) {
	if err = n.DB(ctx).Delete(&MessageDelivery{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrMessageDeliveryDelete, err.Error())
	}
	return
}

// GetById ...
func (n *MessageDeliveries) GetById(ctx context.Context, id int64) (msg *MessageDelivery, err error) {

	msg = &MessageDelivery{}
	err = n.DB(ctx).Model(msg).
		Where("id = ?", id).
		Preload("Message").
		First(&msg).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrMessageDeliveryNotFound, fmt.Sprintf("id \"%d\"", id))
			return
		}
		err = errors.Wrap(apperr.ErrMessageDeliveryGet, err.Error())
	}
	return
}

// MessageDeliveryQuery ...
type MessageDeliveryQuery struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Types     []string   `json:"triggers"`
}
