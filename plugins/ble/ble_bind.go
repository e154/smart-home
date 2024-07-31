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

type Response struct {
	Response []byte `json:"response"`
	Error    string `json:"error"`
}

func GetWriteBind(actor *Actor) func(char string, payload []byte, withResponse bool) Response {
	return func(char string, payload []byte, withResponse bool) Response {
		response, err := actor.ble.Write(char, payload, withResponse)
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{Response: response}
	}
}

func GetReadBind(actor *Actor) func(char string) Response {
	return func(char string) Response {
		response, err := actor.ble.Read(char)
		if err != nil {
			return Response{Error: err.Error()}
		}
		return Response{Response: response}
	}
}
