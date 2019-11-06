package models

import "time"

// swagger:model
type MessageDelivery struct {
	Id                 int64     `json:"id"`
	Message            *Message  `json:"message"`
	MessageId          int64     `json:"message_id"`
	Address            string    `json:"address"`
	Status             string    `json:"status"`
	ErrorMessageStatus *string   `json:"error_message_status"`
	ErrorMessageBody   *string   `json:"error_message_body"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
