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

import "tinygo.org/x/bluetooth"

type WriteGattCharResult struct {
	N     int    `json:"n"`
	Error string `json:"error"`
}

func GetWriteGattCharBind(actor *Actor) func(char string, payload []byte) WriteGattCharResult {
	return func(c string, payload []byte) WriteGattCharResult {
		address, err := bluetooth.ParseUUID(actor.Setts[AttrAddress].String())
		if err != nil {
			return WriteGattCharResult{Error: err.Error()}
		}
		char, err := bluetooth.ParseUUID(c)
		if err != nil {
			return WriteGattCharResult{Error: err.Error()}
		}
		n, err := actor.ble.Write(address, char, payload)
		if err != nil {
			return WriteGattCharResult{Error: err.Error()}
		}
		return WriteGattCharResult{N: n}
	}
}
