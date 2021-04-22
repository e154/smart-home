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
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type IMessageDelivery interface {
	Add(msg *m.MessageDelivery) (id int64, err error)
	SetStatus(msg *m.MessageDelivery) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.MessageDelivery, total int64, err error)
	GetAllUncompleted(limit, offset int64) (list []*m.MessageDelivery, total int64, err error)
	Delete(id int64) (err error)
	GetById(id int64) (ver *m.MessageDelivery, err error)
	fromDb(dbVer *db.MessageDelivery) (ver *m.MessageDelivery)
	toDb(ver *m.MessageDelivery) (dbVer *db.MessageDelivery)
}

// MessageDelivery ...
type MessageDelivery struct {
	table *db.MessageDeliveries
	db    *gorm.DB
}

// GetMessageDeliveryAdaptor ...
func GetMessageDeliveryAdaptor(d *gorm.DB) IMessageDelivery {
	return &MessageDelivery{
		table: &db.MessageDeliveries{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MessageDelivery) Add(msg *m.MessageDelivery) (id int64, err error) {
	id, err = n.table.Add(n.toDb(msg))
	return
}

// SetStatus ...
func (n *MessageDelivery) SetStatus(msg *m.MessageDelivery) (err error) {
	err = n.table.SetStatus(n.toDb(msg))
	return
}

// List ...
func (n *MessageDelivery) List(limit, offset int64, orderBy, sort string) (list []*m.MessageDelivery, total int64, err error) {
	var dbList []*db.MessageDelivery
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MessageDelivery, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// GetAllUncompleted ...
func (n *MessageDelivery) GetAllUncompleted(limit, offset int64) (list []*m.MessageDelivery, total int64, err error) {
	var dbList []*db.MessageDelivery
	if dbList, total, err = n.table.GetAllUncompleted(limit, offset); err != nil {
		return
	}

	list = make([]*m.MessageDelivery, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// Delete ...
func (n *MessageDelivery) Delete(id int64) (err error) {
	err = n.table.Delete(id)
	return
}

// GetById ...
func (n *MessageDelivery) GetById(id int64) (ver *m.MessageDelivery, err error) {

	var dbVer *db.MessageDelivery
	if dbVer, err = n.table.GetById(id); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MessageDelivery) fromDb(dbVer *db.MessageDelivery) (ver *m.MessageDelivery) {

	ver = &m.MessageDelivery{
		Id:                 dbVer.Id,
		MessageId:          dbVer.MessageId,
		Address:            dbVer.Address,
		Status:             m.MessageStatus(dbVer.Status),
		ErrorMessageStatus: dbVer.ErrorMessageStatus,
		ErrorMessageBody:   dbVer.ErrorMessageBody,
		CreatedAt:          dbVer.CreatedAt,
		UpdatedAt:          dbVer.UpdatedAt,
	}

	if dbVer.Message != nil {
		messageAdaptor := GetMessageAdaptor(n.db)
		ver.Message = messageAdaptor.fromDb(dbVer.Message)
	}

	return
}

func (n *MessageDelivery) toDb(ver *m.MessageDelivery) (dbVer *db.MessageDelivery) {

	dbVer = &db.MessageDelivery{
		Id:                 ver.Id,
		MessageId:          ver.MessageId,
		Address:            ver.Address,
		Status:             string(ver.Status),
		ErrorMessageStatus: ver.ErrorMessageStatus,
		ErrorMessageBody:   ver.ErrorMessageBody,
		CreatedAt:          ver.CreatedAt,
		UpdatedAt:          ver.UpdatedAt,
	}

	return
}
