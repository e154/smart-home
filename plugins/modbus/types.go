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

package modbus

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	Name         = "modbus"
	EntityModbus = common.EntityType("modbus")
)

const (
	AttrBaud     = "baud"
	AttrDevice   = "device"
	AttrTimeout  = "timeout"
	AttrStopBits = "stop_bits"
	AttrSleep    = "sleep"
)

func NewAttr() m.EntityAttributes {
	return m.EntityAttributes{
		AttrBaud: {
			Name: AttrBaud,
			Type: common.EntityAttributeInt,
		},
		AttrDevice: {
			Name: AttrDevice,
			Type: common.EntityAttributeInt,
		},
		AttrTimeout: {
			Name: AttrTimeout,
			Type: common.EntityAttributeInt,
		},
		AttrStopBits: {
			Name: AttrStopBits,
			Type: common.EntityAttributeInt,
		},
		AttrSleep: {
			Name: AttrSleep,
			Type: common.EntityAttributeInt,
		},
	}
}
