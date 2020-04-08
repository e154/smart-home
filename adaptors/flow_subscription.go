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

// FlowSubscription ...
type FlowSubscription struct {
	db    *gorm.DB
	table *db.FlowSubscriptions
}

// GetFlowSubscriptionAdaptor ...
func GetFlowSubscriptionAdaptor(Db *gorm.DB) *FlowSubscription {
	return &FlowSubscription{
		db:    Db,
		table: db.NewFlowSubscriptions(Db),
	}
}

// Add ...
func (f *FlowSubscription) Add(sub *m.FlowSubscription) (err error) {
	err = f.table.Add(f.toDb(sub))
	return
}

// Remove ...
func (f *FlowSubscription) Remove(ids []int64) (err error) {
	err = f.table.Delete(ids)
	return
}

func (f *FlowSubscription) fromDb(dbVer *db.FlowSubscription) (ver *m.FlowSubscription) {

	ver = &m.FlowSubscription{
		Id:     dbVer.Id,
		FlowId: dbVer.FlowId,
		Topic:  dbVer.Topic,
	}

	return
}

func (f *FlowSubscription) toDb(ver *m.FlowSubscription) (dbVer *db.FlowSubscription) {

	dbVer = &db.FlowSubscription{
		Id:     ver.Id,
		FlowId: ver.FlowId,
		Topic:  ver.Topic,
	}

	return
}
