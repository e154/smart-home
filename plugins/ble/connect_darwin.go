// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package ble

import (
	"time"

	"tinygo.org/x/bluetooth"
)

func connectBluetooth(address bluetooth.UUID, timeout int64) (*bluetooth.Device, error) {

	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		if err.Error() != "already calling Enable function" {
			return nil, err
		}
	}

	if timeout == 0 {
		timeout = 1
	}

	device, err := adapter.Connect(bluetooth.Address{UUID: address}, bluetooth.ConnectionParams{
		ConnectionTimeout: bluetooth.NewDuration(time.Second * time.Duration(timeout)),
		Timeout:           bluetooth.NewDuration(time.Second * time.Duration(timeout)),
	})
	if err != nil {
		return nil, err
	}

	//log.Infof("connected: %s", device.Address.String())

	return &device, nil
}
