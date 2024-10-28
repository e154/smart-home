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
	"time"

	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

// EntityStorageRepo ...
type EntityStorageRepo interface {
	Add(ctx context.Context, ver *m.EntityStorage) (id int64, err error)
	GetLastByEntityId(ctx context.Context, entityId common.EntityId) (ver *m.EntityStorage, err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string,
		entityIds []common.EntityId,
		startDate, endDate *time.Time) (list []*m.EntityStorage, total int64, err error)
	GetLastThreeById(ctx context.Context, entityId common.EntityId, id int64) (list []*m.EntityStorage, err error)
	DeleteOldest(ctx context.Context, days int) (err error)
}
