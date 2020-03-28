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

package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Nodes struct {
	Db *gorm.DB
}

type Node struct {
	Id                int64 `gorm:"primary_key"`
	Name              string
	Status            string
	Description       string
	Login             string
	EncryptedPassword string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (d *Node) TableName() string {
	return "nodes"
}

func (n Nodes) Add(node *Node) (id int64, err error) {
	if err = n.Db.Create(&node).Error; err != nil {
		return
	}
	id = node.Id
	return
}

func (n Nodes) GetAllEnabled() (list []*Node, err error) {
	list = make([]*Node, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	return
}

func (n Nodes) GetById(nodeId int64) (node *Node, err error) {
	node = &Node{Id: nodeId}
	err = n.Db.First(&node).Error
	return
}

func (n Nodes) Update(m *Node) (err error) {
	q := map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"login":       m.Login,
	}
	if m.EncryptedPassword != "" {
		q["encrypted_password"] = m.EncryptedPassword
	}
	err = n.Db.Model(&Node{Id: m.Id}).Updates(q).Error
	return
}

func (n Nodes) Delete(nodeId int64) (err error) {
	err = n.Db.Delete(&Node{Id: nodeId}).Error
	return
}

func (n *Nodes) List(limit, offset int64, orderBy, sort string) (list []*Node, total int64, err error) {

	if err = n.Db.Model(Node{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Node, 0)
	q := n.Db.Model(&Node{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	return
}

func (n *Nodes) Search(query string, limit, offset int) (list []*Node, total int64, err error) {

	q := n.Db.Model(&Node{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Node, 0)
	err = q.Find(&list).Error

	return
}

func (n *Nodes) GetByLogin(login string) (node *Node, err error) {

	node = &Node{}
	err = n.Db.Model(node).
		Where("login = ?", login).
		First(&node).
		Error

	return
}

func (n *Nodes) GetByName(name string) (node *Node, err error) {

	node = &Node{}
	err = n.Db.Model(node).
		Where("name = ?", name).
		First(&node).
		Error

	return
}
