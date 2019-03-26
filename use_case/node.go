package use_case

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

type NodeCommand struct {
	*CommonCommand
}

func NewNodeCommand(common *CommonCommand) *NodeCommand {
	return &NodeCommand{
		CommonCommand: common,
	}
}

func (n *NodeCommand) Add(params *m.Node) (result *m.Node, errs []*validation.Error, err error) {

	_, errs = node.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.Node.Add(node); err != nil {
		return
	}

	if node, err = n.adaptors.Node.GetById(id); err != nil {
		return
	}

	_, err = n.core.AddNode(node)

	return
}

func (n *NodeCommand) GetById(nodeId int64) (result *m.Node, err error) {

	result, err = n.adaptors.Node.GetById(nodeId)

	return
}

func (n *NodeCommand) Update(params *m.Node) (result *m.Node, errs []*validation.Error, err error) {

	var node *m.Node
	if node, err = n.adaptors.Node.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&node, &params, common.JsonEngine)

	// validation
	_, errs = node.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Node.Update(node); err != nil {
		return
	}

	if node, err = n.adaptors.Node.GetById(node.Id); err != nil {
		return
	}

	// reload node
	err = n.core.ReloadNode(node)

	return
}

func (n *NodeCommand) GetList(limit, offset int64, order, sortBy string) (result []*models.Node, total int64, err error) {

	result, total, err = n.adaptors.Node.List(limit, offset, order, sortBy)

	return
}

func (n *NodeCommand) Delete(nodeId int64) (err error) {

	if nodeId == 0 {
		err = errors.New("node id is null")
		return
	}

	var node *m.Node
	if node, err = n.adaptors.Node.GetById(nodeId); err != nil {
		return
	}

	// remove node from process
	if err = n.core.RemoveNode(node); err != nil {
		return
	}

	err = n.adaptors.Node.Delete(node.Id)

	return
}

func (n *NodeCommand) Search(query string, limit, offset int) (result []*m.Node, total int64, err error) {

	result, total, err = n.adaptors.Node.Search(query, limit, offset)

	return
}
