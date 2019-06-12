package models

import "time"

type DeviceActionScript struct {
	Id int64 `json:"id"`
}
type DeviceActionDevice struct {
	Id int64 `json:"id"`
}

// swagger:model
type NewDeviceAction struct {
	Name        string              `json:"name" valid:"MaxSize(254);Required"`
	Description string              `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
}

// swagger:model
type UpdateDeviceAction struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name" valid:"MaxSize(254);Required"`
	Description string              `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
}

// swagger:model
type DeviceAction struct {
	Id          int64               `json:"id"`
	Name        string              `json:"name" valid:"MaxSize(254);Required"`
	Description string              `json:"description"`
	Device      *DeviceActionDevice `json:"device"`
	Script      *DeviceActionScript `json:"script"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}
