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

package models

// NewDeviceNode ...
type NewDeviceNode struct {
	Id int64 `json:"id"`
}

// swagger:model
type NewDevice struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Device      *ParentDevice    `json:"device"`
	Type        string           `json:"type"`
	Node        *NewDeviceNode   `json:"node"`
	Properties  DeviceProperties `json:"properties"`
}

// swagger:model
type UpdateDevice struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Device      *ParentDevice    `json:"device"`
	Type        string           `json:"type"`
	Node        *NewDeviceNode   `json:"node"`
	Properties  DeviceProperties `json:"properties"`
}

// ParentDevice ...
type ParentDevice struct {
	Id int64 `json:"id"`
}

// swagger:model
type Device struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//Node        *Node            `json:"node"`
	Properties DeviceProperties `json:"properties"`
	Type       string           `json:"type"`
	Status     string           `json:"status"`
	IsGroup    bool             `json:"is_group"`
	Actions    []DeviceAction   `json:"actions"`
	States     []DeviceState    `json:"states"`
	Device     *ParentDevice    `json:"device"`
	DeviceId   *int64           `json:"device_id"`
}

// swagger:model
type DeviceShort struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}
