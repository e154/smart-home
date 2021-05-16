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

package models

import "time"

// MessageType ...
type MessageType string

const (
	// MessageTypeSMS ...
	MessageTypeSMS = MessageType("sms")
	// MessageTypeEmail ...
	MessageTypeEmail = MessageType("email")
	// MessageTypeSlack ...
	MessageTypeSlack = MessageType("slack")
	// MessageTypeUiNotify ...
	MessageTypeUiNotify = MessageType("ui_notify")
	// MessageTypeTelegramNotify ...
	MessageTypeTelegramNotify = MessageType("telegram_notify")
)

// Message ...
type Message struct {
	Id           int64       `json:"id"`
	Type         MessageType `json:"type"`
	EmailFrom    *string     `json:"email_from"`
	EmailSubject *string     `json:"email_subject"`
	EmailBody    *string     `json:"email_body"`
	SmsText      *string     `json:"sms_text"`
	UiText       *string     `json:"ui_text"`
	TelegramText *string     `json:"telegram_text"`
	SlackText    *string     `json:"slack_text"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

// NewNotifrMessage ...
type NewNotifrMessage struct {
	Type         string                 `json:"type"`
	BodyType     string                 `json:"body_type"`
	EmailFrom    *string                `json:"email_from"`
	EmailSubject *string                `json:"email_subject"`
	EmailBody    *string                `json:"email_body"`
	Template     *string                `json:"template"`
	SmsText      *string                `json:"sms_text"`
	SlackText    *string                `json:"slack_text"`
	TelegramText *string                `json:"telegram_text"`
	Params       map[string]interface{} `json:"params"`
	Address      string                 `json:"address"`
}
