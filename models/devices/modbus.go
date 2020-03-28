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

package devices

import (
	. "github.com/e154/smart-home/common"
)

const (
	DevTypeModbusRtu = DeviceType("modbus_rtu")
	DevTypeModbusTcp = DeviceType("modbus_tcp")

	// bit access
	ReadCoils          = "ReadCoils"
	ReadDiscreteInputs = "ReadDiscreteInputs"
	WriteSingleCoil    = "WriteSingleCoil"
	WriteMultipleCoils = "WriteMultipleCoils"

	// 16-bit access
	ReadInputRegisters         = "ReadInputRegisters"
	ReadHoldingRegisters       = "ReadHoldingRegisters"
	ReadWriteMultipleRegisters = "ReadWriteMultipleRegisters"
	WriteSingleRegister        = "WriteSingleRegister"
	WriteMultipleRegisters     = "WriteMultipleRegisters"
)

type DevModBusRtuConfig struct {
	Validation
	SlaveId  int    `json:"slave_id" mapstructure:"slave_id"`   // 1-32
	Baud     int    `json:"baud"`                               // 9600, 19200, ...
	DataBits int    `json:"data_bits" mapstructure:"data_bits"` // 5-9
	StopBits int    `json:"stop_bits" mapstructure:"stop_bits"` // 1, 2
	Parity   string `json:"parity"`                             // none, odd, even
	Timeout  int    `json:"timeout"`                            // milliseconds
}

type DevModBusTcpConfig struct {
	Validation
	SlaveId     int    `json:"slave_id" mapstructure:"slave_id"`
	AddressPort string `json:"address_port" mapstructure:"address_port"`
}

type DevModBusRequest struct {
	Function string   `json:"function"`
	Address  uint16   `json:"address"`
	Count    uint16   `json:"count"`
	Command  []uint16 `json:"command"`
}

// params:
// result
// error
// time
type DevModBusResponse struct {
	BaseResponse
	Result []uint16 `json:"result"`
}

// Javascript Binding
//
// ModBus(func, address, count, command)
//
func NewModBusBind(f string, address, count uint16, command []uint16) ModBusBind {
	return ModBusBind{
		F:       f,
		Address: address,
		Count:   count,
		Command: command,
	}
}

type ModBusBind struct {
	F              string
	Address, Count uint16
	Command        []uint16
}
