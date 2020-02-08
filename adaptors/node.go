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

package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
)

type Node struct {
	table *db.Nodes
	db    *gorm.DB
}

func GetNodeAdaptor(d *gorm.DB) *Node {
	return &Node{
		table: &db.Nodes{Db: d},
		db:    d,
	}
}

func (n *Node) Add(node *m.Node) (id int64, err error) {

	var dbNode *db.Node
	dbNode, err = n.toDb(node)
	if id, err = n.table.Add(dbNode); err != nil {
		return
	}

	return
}

func (n *Node) GetAllEnabled() (list []*m.Node, err error) {

	var dbList []*db.Node
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Node, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (n *Node) GetById(nodeId int64) (node *m.Node, err error) {

	var dbNode *db.Node
	if dbNode, err = n.table.GetById(nodeId); err != nil {
		return
	}

	node = n.fromDb(dbNode)

	return
}

func (n *Node) Update(node *m.Node) (err error) {

	var dbNode *db.Node
	dbNode, err = n.toDb(node)
	err = n.table.Update(dbNode)
	return
}

func (n *Node) Delete(nodeId int64) (err error) {
	err = n.table.Delete(nodeId)
	return
}

func (n *Node) List(limit, offset int64, orderBy, sort string) (list []*m.Node, total int64, err error) {
	var dbList []*db.Node
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Node, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (n *Node) Search(query string, limit, offset int) (list []*m.Node, total int64, err error) {
	var dbList []*db.Node
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Node, 0)
	for _, dbNode := range dbList {
		node := n.fromDb(dbNode)
		list = append(list, node)
	}

	return
}

func (a *Node) GetByLogin(login string) (ver *m.Node, err error) {

	var dbVer *db.Node
	if dbVer, err = a.table.GetByLogin(login); err != nil {
		return
	}

	ver = a.fromDb(dbVer)

	return
}

func (n *Node) fromDb(dbNode *db.Node) (node *m.Node) {
	node = &m.Node{
		Id:                dbNode.Id,
		Name:              dbNode.Name,
		Status:            dbNode.Status,
		Description:       dbNode.Description,
		Login:             dbNode.Login,
		EncryptedPassword: dbNode.EncryptedPassword,
		CreatedAt:         dbNode.CreatedAt,
		UpdatedAt:         dbNode.UpdatedAt,
	}

	return
}

func (n *Node) toDb(node *m.Node) (dbNode *db.Node, err error) {
	dbNode = &db.Node{
		Id:          node.Id,
		Name:        node.Name,
		Status:      node.Status,
		Description: node.Description,
		Login:       node.Login,
	}

	if node.Password != "" {
		if dbNode.EncryptedPassword, err = common.HashPassword(node.Password); err != nil {
			return
		}
	}

	return
}
