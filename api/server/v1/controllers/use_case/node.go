package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"errors"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
	"github.com/e154/smart-home/common"
)

func AddNode(params *models.NewNode, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	node := &m.Node{}
	common.Copy(&node, &params, common.JsonEngine)

	// validation
	ok, errs = node.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.Node.Add(node); err != nil {
		return
	}

	node.Id = id

	// add node
	_, err = core.AddNode(node)

	return
}

func GetNodeById(nodeId int64, adaptors *adaptors.Adaptors) (result *models.NodeModel, err error) {

	var node *m.Node
	if node, err = adaptors.Node.GetById(nodeId); err != nil {
		return
	}

	result = &models.NodeModel{}
	err = common.Copy(&result, &node, common.JsonEngine)

	return
}

func UpdateNode(nodeParams *models.UpdateNodeModel, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	var node *m.Node
	if node, err = adaptors.Node.GetById(nodeParams.Id); err != nil {
		return
	}

	copier.Copy(&node, &nodeParams)

	// validation
	ok, errs = node.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.Node.Update(node); err != nil {
		return
	}

	// reload node
	err = core.ReloadNode(node)

	return
}

func GetNodeList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.NodeModel, total int64, err error) {

	var list []*m.Node
	if list, total, err = adaptors.Node.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.NodeModel, 0)
	err = common.Copy(&result, &list, common.JsonEngine)

	return
}

func DeleteNodeById(nodeId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if nodeId == 0 {
		err = errors.New("node id is null")
		return
	}

	var node *m.Node
	if node, err = adaptors.Node.GetById(nodeId); err != nil {
		return
	}

	// remove node from process
	if err = core.RemoveNode(node); err != nil {
		return
	}

	err = adaptors.Node.Delete(node.Id)

	return
}

func SearchNode(query string, limit, offset int, adaptors *adaptors.Adaptors) (nodes []*m.Node, total int64, err error) {

	nodes, total, err = adaptors.Node.Search(query, limit, offset)

	return
}
