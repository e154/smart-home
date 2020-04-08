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

// NewMessage ...
func NewMessage() (m *Message) {
	m = &Message{
		storage: NewStorage(),
	}
	return
}

// Message ...
type Message struct {
	Error     string
	storage   Storage
	Success   bool
	Direction bool
	Mqtt      bool
}

func (m *Message) clearError() {
	m.Error = ""
}

// SetError ...
func (m *Message) SetError(err string) {
	m.Error = err
}

// Setdir ...
func (m *Message) Setdir(d bool) {
	m.Direction = d
}

// Ok ...
func (m *Message) Ok() {
	m.Success = true
}

// Clear ...
func (m *Message) Clear() {
	m.storage.pull = make(map[string]interface{})
	m.Error = ""
}

// Copy ...
func (m *Message) Copy() (msg *Message) {
	msg = NewMessage()
	for k, v := range m.storage.pull {
		msg.storage.SetVar(k, v)
	}
	return
}

// GetVar ...
func (m *Message) GetVar(key string) (value interface{}) {
	return m.storage.GetVar(key)
}

// SetVar ...
func (m *Message) SetVar(key string, value interface{}) {
	m.storage.SetVar(key, value)
}

// Update ...
func (m *Message) Update(newMsg *Message) {
	m.Error = newMsg.Error
	m.Success = newMsg.Success
	m.Direction = newMsg.Direction
	m.Mqtt = newMsg.Mqtt
	m.storage.copy(newMsg.storage.pull)
}
