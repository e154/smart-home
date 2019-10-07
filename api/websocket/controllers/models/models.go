package models

import (
	//m "github.com/e154/smart-home/models"
)

type DeviceStateStatus struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
	SystemName  string `json:"system_name" valid:"MaxSize(254);Required"`
	DeviceId    int64  `json:"device_id" valid:"Required"`
}

type DeviceState struct {
	Id          int64          `json:"id"`
	Status      *DeviceStateStatus `json:"status"`
	Options     interface{}    `json:"options"`
	ElementName string         `json:"element_name"`
}
