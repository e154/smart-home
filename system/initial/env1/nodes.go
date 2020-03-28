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

package env1

import (
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type NodeManager struct {
	adaptors *adaptors.Adaptors
}

func NewNodeManager(adaptors *adaptors.Adaptors) *NodeManager {
	return &NodeManager{
		adaptors: adaptors,
	}
}

func (n NodeManager) addNode(name, login, pass string) (node *m.Node) {

	var err error
	if node, err = n.adaptors.Node.GetByLogin(name); err == nil || node != nil {
		return
	}

	node = &m.Node{
		Name:     name,
		Status:   "enabled",
		Login:    login,
		Password: pass,
	}

	ok, _ := node.Valid()
	So(ok, ShouldEqual, true)

	node.Id, err = n.adaptors.Node.Add(node)
	So(err, ShouldBeNil)

	return
}

func (n NodeManager) Create() (node1, node2 *m.Node) {

	node1 = n.addNode("node1", "node1", "node1")
	node2 = n.addNode("node2", "node2", "node2")

	return
}
