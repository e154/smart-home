// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package modbus_rtu

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/node"
)

const (
	// Name ...
	Name = "modbus_rtu"
	// FuncEntityAction ...
	FuncEntityAction = "entityAction"
	// DeviceTypeModbusRtu ...
	DeviceTypeModbusRtu = node.DeviceType("modbus_rtu")

	Version = "0.0.1"
)

const (
	// AttrSlaveId ...
	AttrSlaveId = "slave_id"
	// AttrBaud ...
	AttrBaud = "baud"
	// AttrDataBits ...
	AttrDataBits = "data_bits"
	// AttrTimeout ...
	AttrTimeout = "timeout"
	// AttrStopBits ...
	AttrStopBits = "stop_bits"
	// AttrSleep ...
	AttrSleep = "sleep"
	// AttrParity ...
	AttrParity = "parity"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return nil
}

// NewSettings ...
func NewSettings() m.Attributes {
	return m.Attributes{
		AttrBaud: {
			Name: AttrBaud,
			Type: common.AttributeInt,
		},
		AttrSlaveId: {
			Name: AttrSlaveId,
			Type: common.AttributeInt,
		},
		AttrDataBits: {
			Name: AttrDataBits,
			Type: common.AttributeInt,
		},
		AttrTimeout: {
			Name: AttrTimeout,
			Type: common.AttributeInt,
		},
		AttrStopBits: {
			Name: AttrStopBits,
			Type: common.AttributeInt,
		},
		AttrSleep: {
			Name: AttrSleep,
			Type: common.AttributeInt,
		},
		AttrParity: {
			Name: AttrParity,
			Type: common.AttributeString,
		},
	}
}

// ModBusResponse ...
type ModBusResponse struct {
	Error  string   `json:"error"`
	Time   float64  `json:"time"`
	Result []uint16 `json:"result"`
}

// ModBusCommand ...
type ModBusCommand struct {
	Function string   `json:"function"`
	Address  uint16   `json:"address"`
	Count    uint16   `json:"count"`
	Command  []uint16 `json:"command"`
}
