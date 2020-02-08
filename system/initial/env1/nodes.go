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

func nodes(adaptors *adaptors.Adaptors) (node1, node2 *m.Node) {

	node1 = &m.Node{
		Name:     "node1",
		Status:   "enabled",
		Login:    "node1",
		Password: "node1",
	}
	node2 = &m.Node{
		Name:     "node2",
		Status:   "disabled",
		Login:    "node2",
		Password: "node2",
	}

	ok, _ := node1.Valid()
	So(ok, ShouldEqual, true)
	ok, _ = node2.Valid()
	So(ok, ShouldEqual, true)

	var err error
	node1.Id, err = adaptors.Node.Add(node1)
	So(err, ShouldBeNil)

	node2.Id, err = adaptors.Node.Add(node2)
	So(err, ShouldBeNil)

	return
}
