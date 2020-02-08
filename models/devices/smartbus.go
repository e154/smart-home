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
	DevTypeSmartBus = DeviceType("smartbus")
)

type DevSmartBusConfig struct {
	Validation
	Baud     int `json:"baud" valid:"Required"`
	Device   int `json:"device"`
	Timeout  int `json:"timeout" valid:"Required"`
	StopBits int `json:"stop_bits" valid:"Required" mapstructure:"stop_bits"`
	Sleep    int `json:"sleep"`
}

type DevSmartBusRequest struct {
	Command []byte `json:"command"`
}
type DevSmartBusResponse struct {
	BaseResponse
	Result []byte `json:"result"`
}
