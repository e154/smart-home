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
	"tinygo.org/x/bluetooth"
)

type WriteGattCharResult struct {
	Response []byte `json:"response"`
	Error    string `json:"error"`
}

type ReadGattCharResult struct {
	Response []byte `json:"response"`
	Error    string `json:"error"`
}

type GattSubscribeResult struct {
	Error string `json:"error"`
}

type GattDisconnectResult struct {
	Error string `json:"error"`
}

func GetWriteGattCharBind(actor *Actor) func(char string, payload []byte, withResponse bool) WriteGattCharResult {
	return func(c string, payload []byte, withResponse bool) WriteGattCharResult {
		char, err := bluetooth.ParseUUID(c)
		if err != nil {
			return WriteGattCharResult{Error: err.Error()}
		}

		response, err := actor.ble.Write(char, payload, withResponse)
		if err != nil {
			return WriteGattCharResult{Error: err.Error()}
		}
		return WriteGattCharResult{Response: response}
	}
}

func GetReadGattCharBind(actor *Actor) func(char string) ReadGattCharResult {
	return func(c string) ReadGattCharResult {
		char, err := bluetooth.ParseUUID(c)
		if err != nil {
			return ReadGattCharResult{Error: err.Error()}
		}

		response, err := actor.ble.Read(char)
		if err != nil {
			return ReadGattCharResult{Error: err.Error()}
		}
		return ReadGattCharResult{Response: response}
	}
}

func GetSubscribeGattBind(actor *Actor) func(char string, handler func([]byte)) GattSubscribeResult {
	return func(c string, handler func([]byte)) GattSubscribeResult {
		char, err := bluetooth.ParseUUID(c)
		if err != nil {
			return GattSubscribeResult{Error: err.Error()}
		}

		err = actor.ble.Subscribe(char, handler)
		if err != nil {
			return GattSubscribeResult{Error: err.Error()}
		}
		return GattSubscribeResult{}
	}
}

func GetDisconnectBind(actor *Actor) func() GattDisconnectResult {
	return func() GattDisconnectResult {
		if err := actor.ble.Disconnect(); err != nil {
			return GattDisconnectResult{Error: err.Error()}
		}
		return GattDisconnectResult{}
	}
}
