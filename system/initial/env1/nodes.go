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
