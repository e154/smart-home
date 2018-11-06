package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
)

func AddNode(node *m.Node, adaptors *adaptors.Adaptors, core *core.Core) (id int64, Err *Error) {

	err400 := NewError(400)

	// validation
	_, errs := node.Valid()
	err400.ValidationToErrors(errs)

	if err400.Errors() {
		Err = err400
		return
	}

	var err error
	if id, err = adaptors.Node.Add(node); err != nil {
		Err = NewError(500, err)
		return
	}

	node.Id = id

	// add node
	if err = core.AddNode(node); err != nil {
		Err = NewError(500, err)
	}

	return
}

func GetNodeById(nodeId int64, adaptors *adaptors.Adaptors) (node *m.Node, Err *Error) {

	var err error
	if node, err = adaptors.Node.GetById(nodeId); err != nil {
		Err = NewError(500, err)
	}

	return
}

func UpdateNode(node *m.Node, adaptors *adaptors.Adaptors, core *core.Core) (Err *Error) {

	if err := adaptors.Node.Update(node); err != nil {
		Err = NewError(500, err)
		return
	}

	// reload node
	if err := core.ReloadNode(node); err != nil {
		Err = NewError(500, err)
	}

	return
}

func GetNodeList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Node, total int64, Err *Error) {

	var err error
	items, total, err = adaptors.Node.List(limit, offset, order, sortBy)
	if err != nil {
		Err = NewError(500, err)
	}

	return
}

func DeleteNodeById(nodeId int64, adaptors *adaptors.Adaptors, core *core.Core) (Err *Error) {

	if nodeId == 0 {
		Err = NewError(400, "node id is null")
		return
	}

	node, err := adaptors.Node.GetById(nodeId)
	if err != nil {
		Err = NewError(500, err)
		return
	}

	// remove node from process
	if err := core.RemoveNode(node); err != nil {
		Err = NewError(500, err)
		return
	}

	if err := adaptors.Node.Delete(node.Id); err != nil {
		Err = NewError(500, err)
	}

	return
}

