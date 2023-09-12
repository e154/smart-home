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

package adaptors

import (
	"context"
	"encoding/json"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"gorm.io/gorm"
)

// IMessage ...
type IMessage interface {
	Add(ctx context.Context, msg *m.Message) (id int64, err error)
	fromDb(dbVer *db.Message) (ver *m.Message)
	toDb(ver *m.Message) (dbVer *db.Message)
}

// Message ...
type Message struct {
	IMessage
	table *db.Messages
}

// GetMessageAdaptor ...
func GetMessageAdaptor(d *gorm.DB) IMessage {
	return &Message{
		table: &db.Messages{Db: d},
	}
}

// Add ...
func (n *Message) Add(ctx context.Context, msg *m.Message) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(msg))
	return
}

func (n *Message) fromDb(dbVer *db.Message) (ver *m.Message) {
	ver = &m.Message{
		Id:         dbVer.Id,
		Type:       dbVer.Type,
		CreatedAt:  dbVer.CreatedAt,
		UpdatedAt:  dbVer.UpdatedAt,
		Attributes: m.AttributeValue{},
	}

	if len(dbVer.Payload) > 0 {
		_ = json.Unmarshal(dbVer.Payload, &ver.Attributes)
	}

	return
}

func (n *Message) toDb(ver *m.Message) (dbVer *db.Message) {
	dbVer = &db.Message{
		Id:        ver.Id,
		Type:      ver.Type,
		CreatedAt: ver.CreatedAt,
		UpdatedAt: ver.UpdatedAt,
	}

	if ver.Attributes != nil {
		b, _ := json.Marshal(ver.Attributes)
		_ = dbVer.Payload.UnmarshalJSON(b)
	}

	return
}
