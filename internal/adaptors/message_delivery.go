// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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
	"context"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.MessageDeliveryRepo = (*MessageDelivery)(nil)

// MessageDelivery ...
type MessageDelivery struct {
	table *db.MessageDeliveries
	db    *gorm.DB
}

// GetMessageDeliveryAdaptor ...
func GetMessageDeliveryAdaptor(d *gorm.DB) *MessageDelivery {
	return &MessageDelivery{
		table: &db.MessageDeliveries{&db.Common{Db: d}},
		db:    d,
	}
}

// Add ...
func (n *MessageDelivery) Add(ctx context.Context, msg *m.MessageDelivery) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(msg))
	return
}

// SetStatus ...
func (n *MessageDelivery) SetStatus(ctx context.Context, msg *m.MessageDelivery) (err error) {
	err = n.table.SetStatus(ctx, n.toDb(msg))
	return
}

// List ...
func (n *MessageDelivery) List(ctx context.Context, limit, offset int64, orderBy, sort string, query *m.MessageDeliveryQuery) (list []*m.MessageDelivery, total int64, err error) {
	var dbList []*db.MessageDelivery
	var queryObj *db.MessageDeliveryQuery
	if query != nil {
		queryObj = &db.MessageDeliveryQuery{
			StartDate: query.StartDate,
			EndDate:   query.EndDate,
			Types:     query.Types,
		}
	}
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, queryObj); err != nil {
		return
	}

	list = make([]*m.MessageDelivery, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// GetAllUncompleted ...
func (n *MessageDelivery) GetAllUncompleted(ctx context.Context, limit, offset int64) (list []*m.MessageDelivery, total int64, err error) {
	var dbList []*db.MessageDelivery
	if dbList, total, err = n.table.GetAllUncompleted(ctx, int(limit), int(offset)); err != nil {
		return
	}

	list = make([]*m.MessageDelivery, 0)
	for _, dbVer := range dbList {
		list = append(list, n.fromDb(dbVer))
	}

	return
}

// Delete ...
func (n *MessageDelivery) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
	return
}

// GetById ...
func (n *MessageDelivery) GetById(ctx context.Context, id int64) (ver *m.MessageDelivery, err error) {

	var dbVer *db.MessageDelivery
	if dbVer, err = n.table.GetById(ctx, id); err != nil {
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
		EntityId:           dbVer.EntityId,
		Status:             m.MessageStatus(dbVer.Status),
		ErrorMessageStatus: dbVer.ErrorMessageStatus,
		ErrorMessageBody:   dbVer.ErrorMessageBody,
		CreatedAt:          dbVer.CreatedAt,
		UpdatedAt:          dbVer.UpdatedAt,
	}

	if dbVer.MessageId != 0 {
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
		EntityId:           ver.EntityId,
		Status:             string(ver.Status),
		ErrorMessageStatus: ver.ErrorMessageStatus,
		ErrorMessageBody:   ver.ErrorMessageBody,
		CreatedAt:          ver.CreatedAt,
		UpdatedAt:          ver.UpdatedAt,
	}

	return
}
