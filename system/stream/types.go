// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package stream

import "github.com/e154/smart-home/common/uuid"

// BroadcastClient ...
type BroadcastClient interface {
	Broadcast(query string, message []byte)
}

// IStreamClient ...
type IStreamClient interface {
	Send(id string, query string, body []byte) error
	Notify(t, b string)
}

// Message ...
type Message struct {
	Id      uuid.UUID              `json:"id"`
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload"`
	Forward string                 `json:"forward"`
	Status  string                 `json:"status"`
	Type    string                 `json:"type"`
}
