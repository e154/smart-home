package core

import (
	"github.com/e154/smart-home/api/models"
)

// Javascript Binding
//
// Device
//	.getName()
//	.getDescription()
//	.getAddress()
//
type DeviceBind struct {
	model 		*models.Device
}

func (d *DeviceBind) GetName() string {
	return d.model.Name
}

func (d *DeviceBind) GetDescription() string {
	return d.model.Description
}

func (d *DeviceBind) GetAddress() *int {
	return d.model.Address
}
