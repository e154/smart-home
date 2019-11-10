package models

import "time"

type MessageStatus string

const (
	MessageStatusNew        = MessageStatus("new")
	MessageStatusInProgress = MessageStatus("in_progress")
	MessageStatusSucceed    = MessageStatus("succeed")
	MessageStatusError      = MessageStatus("error")
)

type MessageDelivery struct {
	Id                 int64         `json:"id"`
	Message            *Message      `json:"message"`
	MessageId          int64         `json:"message_id"`
	Address            string        `json:"address"`
	Status             MessageStatus `json:"status"`
	ErrorMessageStatus *string       `json:"error_message_status"`
	ErrorMessageBody   *string       `json:"error_message_body"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}
