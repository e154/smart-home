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

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/system/uuid"
)

// Message ...
type Message struct {
	Id      uuid.UUID              `json:"id"`
	Command string                 `json:"command"`
	Payload map[string]interface{} `json:"payload"`
	Forward string                 `json:"forward"`
	Status  string                 `json:"status"`
	Type    string                 `json:"type"`
}

// NewMessage ...
func NewMessage(b []byte) (message Message, err error) {

	message = Message{}
	err = json.Unmarshal(b, &message)

	return
}

// Pack ...
func (m *Message) Pack() []byte {
	b, _ := json.Marshal(m)
	return b
}

// Response ...
func (m *Message) Response(payload map[string]interface{}) *Message {
	msg := &Message{
		Id:      m.Id,
		Payload: payload,
		Forward: Response,
		Status:  StatusSuccess,
	}
	return msg
}

// Success ...
func (m *Message) Success() *Message {
	msg := &Message{
		Id:      m.Id,
		Payload: map[string]interface{}{},
		Forward: Response,
		Status:  StatusSuccess,
	}
	return msg
}

// Error ...
func (m *Message) Error(err error) *Message {
	msg := &Message{
		Id: m.Id,
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
		Forward: Response,
		Status:  StatusError,
	}
	return msg
}

// IsError ...
func (m *Message) IsError() (err error) {

	if m.Status != StatusError {
		return
	}
	if _, ok := m.Payload["error"]; !ok {
		return
	}
	err = fmt.Errorf("%v", m.Payload["error"])
	return
}
