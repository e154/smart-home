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
// IC.Flow()
//	 .getName()
//	 .getDescription()
//	 .setVar(string, interface)
//	 .getVar(string)
//	 .node()
//
type FlowBind struct {
	flow *Flow
}

func (f *FlowBind) GetName() string {
	return f.flow.Model.Name
}

func (f *FlowBind) GetDescription() string {
	return f.flow.Model.Description
}

func (f *FlowBind) SetVar(key string, value interface{}) {
	f.flow.SetVar(key, value)
}

func (f *FlowBind) GetVar(key string) interface{} {
	return f.flow.GetVar(key)
}

func (f *FlowBind) Node() *NodeBind {
	if f.flow.Node == nil {
		return nil
	}
	return &NodeBind{node: f.flow.Node}
}
