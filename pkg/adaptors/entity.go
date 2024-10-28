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

	"github.com/e154/smart-home/pkg/common"
	m "github.com/e154/smart-home/pkg/models"
)

// EntityRepo ...
type EntityRepo interface {
	Add(ctx context.Context, ver *m.Entity) (err error)
	GetById(ctx context.Context, id common.EntityId, preloadMetric ...bool) (ver *m.Entity, err error)
	GetByIds(ctx context.Context, ids []common.EntityId, preloadMetric ...bool) (ver []*m.Entity, err error)
	GetByIdsSimple(ctx context.Context, ids []common.EntityId) (list []*m.Entity, err error)
	Delete(ctx context.Context, id common.EntityId) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string, autoLoad bool, query, plugin *string,
		areaId *int64) (list []*m.Entity, total int64, err error)
	ListPlain(ctx context.Context, limit, offset int64, orderBy, sort string, autoLoad bool, query, plugin *string,
		areaId *int64, tags *[]string) (list []*m.Entity, total int64, err error)
	GetByType(ctx context.Context, t string, limit, offset int64) (list []*m.Entity, err error)
	Update(ctx context.Context, ver *m.Entity) (err error)
	Search(ctx context.Context, query string, limit, offset int64) (list []*m.Entity, total int64, err error)
	UpdateAutoload(ctx context.Context, entityId common.EntityId, autoLoad bool) (err error)
	Statistic(ctx context.Context) (statistic *m.EntitiesStatistic, err error)
	DeleteScripts(ctx context.Context, entityID common.EntityId) (err error)
	DeleteTags(ctx context.Context, entityID common.EntityId) (err error)
}
