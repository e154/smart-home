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
	"encoding/json"

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.MessageRepo = (*Message)(nil)

// Message ...
type Message struct {
	table *db.Messages
}

// GetMessageAdaptor ...
func GetMessageAdaptor(d *gorm.DB) *Message {
	return &Message{
		table: &db.Messages{&db.Common{Db: d}},
	}
}

// Add ...
func (n *Message) Add(ctx context.Context, msg *models.Message) (id int64, err error) {
	id, err = n.table.Add(ctx, n.toDb(msg))
	return
}

func (n *Message) fromDb(dbVer *db.Message) (ver *models.Message) {
	ver = &models.Message{
		Id:         dbVer.Id,
		Type:       dbVer.Type,
		CreatedAt:  dbVer.CreatedAt,
		UpdatedAt:  dbVer.UpdatedAt,
		Attributes: models.AttributeValue{},
	}

	if len(dbVer.Payload) > 0 {
		_ = json.Unmarshal(dbVer.Payload, &ver.Attributes)
	}

	return
}

func (n *Message) toDb(ver *models.Message) (dbVer *db.Message) {
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
