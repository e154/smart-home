package models

import "time"

type MessageType string

const (
	MessageTypeSMS            = MessageType("sms")
	MessageTypeEmail          = MessageType("email")
	MessageTypeSlack          = MessageType("slack")
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
	SlackText    *string     `json:"slack_text"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

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
