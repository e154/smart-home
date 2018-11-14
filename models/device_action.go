package models

import "time"

type DeviceAction struct {
	Id          int64     `json:"id"`
	Device      *Device   `json:"device"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Script      *Script   `json:"script"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
