package models

import "time"

type MessageType string

const (
	MessageTypeSMS            = MessageType("sms")
	MessageTypeEmail          = MessageType("email")
	MessageTypeUiNotify       = MessageType("ui_notify")
	MessageTypeTelegramNotify = MessageType("telegram_notify")
)

type Message struct {
	Id           int64       `json:"id"`
	Type         MessageType `json:"type"`
	EmailFrom    *string     `json:"email_from"`
	EmailSubject *string     `json:"email_subject"`
	EmailBody    *string     `json:"email_body"`
	SmsText      *string     `json:"sms_text"`
	UiText       *string     `json:"ui_text"`
	TelegramText *string     `json:"telegram_text"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
