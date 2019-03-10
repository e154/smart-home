package models

import "time"

type MapDevice struct {
	Id         int64              `json:"id"`
	SystemName string             `json:"system_name" valid:"Required"`
	Device     *DeviceModel       `json:"device"`
	DeviceId   int64              `json:"device_id" valid:"Required"`
	Image      *Image             `json:"image"`
	ImageId    int64              `json:"image_id"`
	Actions    []*MapDeviceAction `json:"actions"`
	States     []*MapDeviceState  `json:"states"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
