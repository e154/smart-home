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

package common

import "strings"

// EntityId ...
type EntityId string

// NewEntityId ...
func NewEntityId(s string) *EntityId {
	e := EntityId(s)
	return &e
}

// NewEntityIdFromPtr ...
func NewEntityIdFromPtr(s *string) *EntityId {
	if s == nil {
		return nil
	}
	e := EntityId(*s)
	return &e
}

// Name ...
func (e EntityId) Name() string {
	arr := strings.Split(string(e), ".")
	if len(arr) > 1 {
		return arr[1]
	}
	return string(e)
}

// PluginName ...
func (e EntityId) PluginName() string {
	arr := strings.Split(string(e), ".")
	if len(arr) > 1 {
		return arr[0]
	}
	return string(e)
}

// String ...
func (e *EntityId) String() string {
	if e == nil {
		return ""
	} else {
		return string(*e)
	}
}

// StringPtr ...
func (e *EntityId) StringPtr() *string {
	if e == nil {
		return nil
	}
	r := e.String()
	return &r
}

// Ptr ...
func (e EntityId) Ptr() *EntityId {
	return &e
}
