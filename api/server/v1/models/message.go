package models

import "time"

// swagger:model
type Message struct {
	Id           int64     `json:"id"`
	Type         string    `json:"type"`
	EmailFrom    *string   `json:"email_from"`
	EmailSubject *string   `json:"email_subject"`
	EmailBody    *string   `json:"email_body"`
	SmsText      *string   `json:"sms_text"`
	UiText       *string   `json:"ui_text"`
	SlackText    *string   `json:"slack_text"`
	TelegramText *string   `json:"telegram_text"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
