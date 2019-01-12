package models

import "time"

type DeviceStateDevice struct {
	Id int64 `json:"id"`
}
type NewDeviceState struct {
	Description string             `json:"description"`
	SystemName  string             `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *DeviceStateDevice `json:"device" valid:"Required"`
}

type UpdateDeviceState struct {
	Description string             `json:"description"`
	SystemName  string             `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *DeviceStateDevice `json:"device" valid:"Required"`
}

type DeviceState struct {
	Id          int64              `json:"id"`
	Description string             `json:"description"`
	SystemName  string             `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *DeviceStateDevice `json:"device" valid:"Required"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}