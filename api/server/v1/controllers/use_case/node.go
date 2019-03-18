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

func AddNode(params *models.NewNode, adaptors *adaptors.Adaptors, core *core.Core) (result *models.Node, errs []*validation.Error, err error) {

	node := &m.Node{}
	common.Copy(&node, &params, common.JsonEngine)

	// validation
	_, errs = node.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = adaptors.Node.Add(node); err != nil {
		return
	}

	if node, err = adaptors.Node.GetById(id); err != nil {
		return
	}

	result = &models.Node{}
	if err = common.Copy(&result, &node, common.JsonEngine); err != nil {
		return
	}

	// add node
	_, err = core.AddNode(node)

	return
}

func GetNodeById(nodeId int64, adaptors *adaptors.Adaptors) (result *models.Node, err error) {

	var node *m.Node
	if node, err = adaptors.Node.GetById(nodeId); err != nil {
		return
	}

	result = &models.Node{}
	err = common.Copy(&result, &node, common.JsonEngine)

	return
}

func UpdateNode(params *models.UpdateNode, adaptors *adaptors.Adaptors, core *core.Core) (result *models.Node, errs []*validation.Error, err error) {

	var node *m.Node
	if node, err = adaptors.Node.GetById(params.Id); err != nil {
		return
	}

	copier.Copy(&node, &params)

	// validation
	_, errs = node.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Node.Update(node); err != nil {
		return
	}

	if node, err = adaptors.Node.GetById(node.Id); err != nil {
		return
	}

	result = &models.Node{}
	if err = common.Copy(&result, &node, common.JsonEngine); err != nil {
		return
	}

	// reload node
	err = core.ReloadNode(node)

	return
}

func GetNodeList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.Node, total int64, err error) {

	var list []*m.Node
	if list, total, err = adaptors.Node.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.Node, 0)
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
