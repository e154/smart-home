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

// Javascript Binding
//
// MapElement
//	.SetState(name)
//	.GetState()
//	.SetOptions(options)
//	.GetOptions()
//	.Story(logLevel, type, description)
//	.PushMetric(name, obj)
//
type MapElementBind struct {
	element *MapElement
}

// SetState ...
func (e *MapElementBind) SetState(name string) {
	e.element.SetState(name)
}

// GetState ...
func (e *MapElementBind) GetState() interface{} {
	return e.element.State
}

// SetOptions ...
func (e *MapElementBind) SetOptions(options interface{}) {
	e.element.SetOptions(options)
}

// GetOptions ...
func (e *MapElementBind) GetOptions() interface{} {
	return e.element.GetOptions()
}

// Story ...
func (e *MapElementBind) Story(logLevel, t, desc string) {
	e.element.CustomHistory(logLevel, t, desc)
}

// PushMetric ...
func (e *MapElementBind) PushMetric(name string, val map[string]interface{}) {
	e.element.PushMetric(name, val)
}
