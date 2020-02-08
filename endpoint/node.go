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

package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type NodeEndpoint struct {
	*CommonEndpoint
}

func NewNodeEndpoint(common *CommonEndpoint) *NodeEndpoint {
	return &NodeEndpoint{
		CommonEndpoint: common,
	}
}

func (n *NodeEndpoint) Add(params *m.Node) (result *m.Node, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.Node.Add(params); err != nil {
		return
	}

	if result, err = n.adaptors.Node.GetById(id); err != nil {
		return
	}

	_, err = n.core.AddNode(result)

	return
}

func (n *NodeEndpoint) GetById(nodeId int64) (result *m.Node, err error) {

	result, err = n.adaptors.Node.GetById(nodeId)

	return
}

func (n *NodeEndpoint) Update(params *m.Node) (result *m.Node, errs []*validation.Error, err error) {

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

func (n *NodeEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.Node, total int64, err error) {

	result, total, err = n.adaptors.Node.List(limit, offset, order, sortBy)

	return
}

func (n *NodeEndpoint) Delete(nodeId int64) (err error) {

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

func (n *NodeEndpoint) Search(query string, limit, offset int) (result []*m.Node, total int64, err error) {

	result, total, err = n.adaptors.Node.Search(query, limit, offset)

	return
}
