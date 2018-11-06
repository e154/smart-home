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
	core.AddNode(node)

	return
}

func GetNodeById() {

}

func GetNodeList() {

}

func DeleteNodeById() {

}

