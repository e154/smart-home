package models

import "time"

type MapDeviceAction struct {
	Id             int64         `json:"id"`
	DeviceAction   *DeviceAction `json:"device_action"`
	DeviceActionId int64         `json:"device_action_id" valid:"Required"`
	MapDeviceId    int64         `json:"map_device_id" valid:"Required"`
	Image          *Image        `json:"image"`
	ImageId        int64         `json:"image_id" valid:"Required"`
	Type           string        `json:"type"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}


