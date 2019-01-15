package core

import (
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
)

// Javascript Binding
//
// Device
//	.getName()
//	.getDescription()
//	.runCommand(command []string)
//	.smartBus(command []byte)
//	.modBus(func string, address, count int64, command []byte)
//
type DeviceBind struct {
	model *m.Device
	node *Node
}

func (d *DeviceBind) GetName() string {
	return d.model.Name
}

func (d *DeviceBind) GetDescription() string {
	return d.model.Description
}

func (d *DeviceBind) RunCommand(name string, args []string) (result *DevCommandResponse) {
	dev := &Device{
		dev: d.model,
		node: d.node,
	}
	result = dev.RunCommand(name, args)

	return
}

func (d *DeviceBind) SmartBus(command []byte) (result *DevSmartBusResponse) {
	dev := &Device{
		dev: d.model,
		node: d.node,
	}
	result = dev.SmartBus(command)
	return
}

func (d *DeviceBind) ModBus(f string, address, count uint16, command []uint16) (result *DevModBusResponse) {
	dev := &Device{
		dev: d.model,
		node: d.node,
	}
	result = dev.ModbusBus(f, address, count, command)
	return
}