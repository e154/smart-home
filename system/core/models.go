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
	"encoding/json"
	"github.com/e154/smart-home/common"
	"time"
)

// NodeMessage ...
type NodeMessage struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Command    json.RawMessage   `json:"command"`
}

// NodeResponse ...
type NodeResponse struct {
	DeviceId   int64             `json:"device_id"`
	DeviceType common.DeviceType `json:"device_type"`
	Properties json.RawMessage   `json:"properties"`
	Response   json.RawMessage   `json:"response"`
	Status     string            `json:"status"`
	Time       float64           `json:"time"`
}

// NodeStatus ...
type NodeStatus string

// NodeStatModel ...
type NodeStatModel struct {
	Status    NodeStatus `json:"status"`
	Thread    int        `json:"thread"`
	Rps       int64      `json:"rps"`
	Min       int64      `json:"min"`
	Max       int64      `json:"max"`
	StartedAt time.Time  `json:"started_at"`
}

// NodeStat ...
type NodeStat struct {
	Status      NodeStatus `json:"status"`
	Thread      int        `json:"thread"`
	Rps         int64      `json:"rps"`
	Min         int64      `json:"min"`
	Max         int64      `json:"max"`
	StartedAt   time.Time  `json:"started_at"`
	LastPing    time.Time  `json:"last_ping"`
	ConnStatus  string     `json:"conn_status"`
	IsConnected bool       `json:"is_connected"`
}
