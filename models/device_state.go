package models

import "time"

type DeviceState struct {
	Id          int64     `json:"id"`
	Device      *Device   `json:"device"`
	Description string    `json:"description"`
	SystemName  string    `json:"system_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
