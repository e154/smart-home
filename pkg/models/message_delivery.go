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

package models

import (
	"time"

	"github.com/e154/smart-home/pkg/common"
)

// MessageStatus ...
type MessageStatus string

const (
	// MessageStatusNew ...
	MessageStatusNew = MessageStatus("new")
	// MessageStatusInProgress ...
	MessageStatusInProgress = MessageStatus("in_progress")
	// MessageStatusSucceed ...
	MessageStatusSucceed = MessageStatus("succeed")
	// MessageStatusError ...
	MessageStatusError = MessageStatus("error")
)

// MessageDeliveryQuery ...
type MessageDeliveryQuery struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Types     []string   `json:"triggers"`
}

// MessageDelivery ...
type MessageDelivery struct {
	Id                 int64            `json:"id"`
	Message            *Message         `json:"message"`
	MessageId          int64            `json:"message_id"`
	Address            string           `json:"address"`
	EntityId           *common.EntityId `json:"entity_id"`
	Status             MessageStatus    `json:"status"`
	ErrorMessageStatus *string          `json:"error_message_status"`
	ErrorMessageBody   *string          `json:"error_message_body"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
}
