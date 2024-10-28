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

	m "github.com/e154/smart-home/pkg/models"
)

// Zigbee2mqttDeviceRepo ...
type Zigbee2mqttDeviceRepo interface {
	Add(ctx context.Context, ver *m.Zigbee2mqttDevice) (err error)
	GetById(ctx context.Context, id string) (ver *m.Zigbee2mqttDevice, err error)
	Update(ctx context.Context, ver *m.Zigbee2mqttDevice) (err error)
	Delete(ctx context.Context, id string) (err error)
	List(ctx context.Context, limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error)
	ListByBridgeId(ctx context.Context, bridgeId, limit, offset int64, orderBy, sort string) (list []*m.Zigbee2mqttDevice, total int64, err error)
	Search(ctx context.Context, query string, limit, offset int64) (list []*m.Zigbee2mqttDevice, total int64, err error)
}
