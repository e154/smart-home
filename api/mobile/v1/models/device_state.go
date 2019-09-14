package models

type DeviceStateDevice struct {
	Id int64 `json:"id"`
}

// swagger:model
type NewDeviceState struct {
	Description string             `json:"description"`
	SystemName  string             `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *DeviceStateDevice `json:"device" valid:"Required"`
}

// swagger:model
type UpdateDeviceState struct {
	Description string             `json:"description"`
	SystemName  string             `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *DeviceStateDevice `json:"device" valid:"Required"`
}

// swagger:model
type DeviceState struct {
	Id          int64              `json:"id"`
	Description string             `json:"description"`
	SystemName  string             `json:"system_name" valid:"MaxSize(254);Required"`
	Device      *DeviceStateDevice `json:"device" valid:"Required"`
}
