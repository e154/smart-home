// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/pkg/errors"
	"time"
)

// Nodes ...
type Nodes struct {
	Db *gorm.DB
}

// Node ...
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

// TableName ...
func (d *Node) TableName() string {
	return "nodes"
}

// Add ...
func (n Nodes) Add(node *Node) (id int64, err error) {
	if err = n.Db.Create(&node).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = node.Id
	return
}

// GetAllEnabled ...
func (n Nodes) GetAllEnabled() (list []*Node, err error) {
	list = make([]*Node, 0)
	err = n.Db.Where("status = ?", "enabled").
		Find(&list).Error
	if err != nil {
		err = errors.Wrap(err, "getAllEnabled failed")
	}
	return
}

// GetById ...
func (n Nodes) GetById(nodeId int64) (node *Node, err error) {
	node = &Node{Id: nodeId}
	if err = n.Db.First(&node).Error; err != nil {
		err = errors.Wrap(err, "getById failed")
	}
	return
}

// Update ...
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
	if err = n.Db.Model(&Node{Id: m.Id}).Updates(q).Error; err != nil {
		err = errors.Wrap(err, "update failed")
	}
	return
}

// Delete ...
func (n Nodes) Delete(nodeId int64) (err error) {
	if err = n.Db.Delete(&Node{Id: nodeId}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
	}
	return
}

// List ...
func (n *Nodes) List(limit, offset int64, orderBy, sort string) (list []*Node, total int64, err error) {

	if err = n.Db.Model(Node{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
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

	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(err, "list failed")
	}

	return
}

// Search ...
func (n *Nodes) Search(query string, limit, offset int) (list []*Node, total int64, err error) {

	q := n.Db.Model(&Node{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Node, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(err, "search failed")
	}

	return
}

// GetByLogin ...
func (n *Nodes) GetByLogin(login string) (node *Node, err error) {

	node = &Node{}
	err = n.Db.Model(node).
		Where("login = ?", login).
		First(&node).
		Error
	if err != nil {
		err = errors.Wrap(err, "getByLogin failed")
	}
	return
}

// GetByName ...
func (n *Nodes) GetByName(name string) (node *Node, err error) {

	node = &Node{}
	err = n.Db.Model(node).
		Where("name = ?", name).
		First(&node).
		Error
	if err != nil {
		err = errors.Wrap(err, "getByName failed")
	}
	return
}
