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
// Node
//	 .Name()
//	 .Status()
//	 .Stat()
//	 .Description()
//	 .IsConnected()
//
type NodeBind struct {
	node *Node
}

//func (n *NodeBind) Send(device *DeviceBind, command []byte) NodeBindResult {
//	return n.node.Send(device.model, command)
//}

func (n *NodeBind) IsConnected() bool {
	return n.node.IsConnected()
}

func (n *NodeBind) Name() string {
	return n.node.Model().Name
}

func (n *NodeBind) Status() string {
	return n.node.Model().Status
}

func (n *NodeBind) Stat() NodeStat {
	return n.node.GetStat()
}

func (n *NodeBind) Description() string {
	return n.node.Model().Description
}
