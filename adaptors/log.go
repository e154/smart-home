// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// ILog ...
type ILog interface {
	Add(ctx context.Context, ver *m.Log) (id int64, err error)
	AddMultiple(ctx context.Context, items []*m.Log) (err error)
	GetById(ctx context.Context, verId int64) (ver *m.Log, err error)
	Delete(ctx context.Context, verId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, queryObj *m.LogQuery) (list []*m.Log, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Log, total int64, err error)
	DeleteOldest(ctx context.Context, days int) (err error)
	fromDb(dbVer *db.Log) (ver *m.Log)
	toDb(ver *m.Log) (dbVer *db.Log)
}

// Log ...
type Log struct {
	ILog
	table *db.Logs
	db    *gorm.DB
}

// GetLogAdaptor ...
func GetLogAdaptor(d *gorm.DB) ILog {
	return &Log{
		table: &db.Logs{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Log) Add(ctx context.Context, ver *m.Log) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(ctx, dbVer); err != nil {
		return
	}

	return
}

// AddMultiple ...
func (n *Log) AddMultiple(ctx context.Context, items []*m.Log) (err error) {

	insertRecords := make([]*db.Log, 0, len(items))
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = n.table.AddMultiple(ctx, insertRecords)

	return
}

// GetById ...
func (n *Log) GetById(ctx context.Context, verId int64) (ver *m.Log, err error) {

	var dbVer *db.Log
	if dbVer, err = n.table.GetById(ctx, verId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Delete ...
func (n *Log) Delete(ctx context.Context, verId int64) (err error) {
	err = n.table.Delete(ctx, verId)
	return
}

// List ...
func (n *Log) List(ctx context.Context, limit, offset int64, orderBy, sort string, queryObj *m.LogQuery) (list []*m.Log, total int64, err error) {

	var dbList []*db.Log
	var dbQueryObj *db.LogQuery

	if queryObj != nil {
		dbQueryObj = &db.LogQuery{
			StartDate: queryObj.StartDate,
			EndDate:   queryObj.EndDate,
			Levels:    queryObj.Levels,
		}
	}

	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort, dbQueryObj); err != nil {
		return
	}

	list = make([]*m.Log, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// Search ...
func (n *Log) Search(ctx context.Context, query string, limit, offset int) (list []*m.Log, total int64, err error) {
	var dbList []*db.Log
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Log, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// DeleteOldest ...
func (n *Log) DeleteOldest(ctx context.Context, days int) (err error) {
	err = n.table.DeleteOldest(ctx, days)
	return
}

func (n *Log) fromDb(dbVer *db.Log) (ver *m.Log) {
	ver = &m.Log{
		Id:        dbVer.Id,
		Body:      dbVer.Body,
		Level:     dbVer.Level,
		Owner:     dbVer.Owner,
		CreatedAt: dbVer.CreatedAt,
	}

	return
}

func (n *Log) toDb(ver *m.Log) (dbVer *db.Log) {
	dbVer = &db.Log{
		Body:  ver.Body,
		Level: ver.Level,
		Owner: ver.Owner,
	}
	return
}
