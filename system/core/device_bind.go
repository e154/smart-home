// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package core

import (
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
)

// Javascript Binding
//
// Device
//	.getName()
//	.getModel()
//	.getDescription()
//	.runCommand(command []string)
//	.smartBus(command []byte)
//	.modBus(func string, address, count int64, command []byte)
//
type DeviceBind struct {
	model *m.Device
	node  *Node
}

func (d *DeviceBind) GetName() string {
	return d.model.Name
}

func (d *DeviceBind) GetModel() *m.Device {
	return d.model
}

func (d *DeviceBind) GetDescription() string {
	return d.model.Description
}

func (d *DeviceBind) RunCommand(name string, args []string) (result *DevCommandResponse) {
	dev := &Device{
		dev:  d.model,
		node: d.node,
	}
	result = dev.RunCommand(name, args)

	return
}

func (d *DeviceBind) SmartBus(command []byte) (result *DevSmartBusResponse) {
	dev := &Device{
		dev:  d.model,
		node: d.node,
	}
	result = dev.SmartBus(command)
	return
}

func (d *DeviceBind) ModBus(f string, address, count uint16, command []uint16) (result *DevModBusResponse) {
	dev := &Device{
		dev:  d.model,
		node: d.node,
	}
	result = dev.ModBus(f, address, count, command)
	return
}
