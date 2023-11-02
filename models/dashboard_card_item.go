// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package models

import (
	"encoding/json"
	"time"

	"github.com/e154/smart-home/common"
)

// DashboardCardItem ...
type DashboardCardItem struct {
	Payload         json.RawMessage  `json:"payload"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	Title           string           `json:"title" validate:"required"`
	Type            string           `json:"type" validate:"required"`
	Id              int64            `json:"id"`
	Weight          int              `json:"weight"`
	DashboardCardId int64            `json:"dashboard_card_id" validate:"required"`
	DashboardCard   *DashboardCard   `json:"dashboard_card"`
	EntityId        *common.EntityId `json:"entity_id"`
	Entity          *Entity          `json:"entity"`
	Enabled         bool             `json:"enabled"`
	Hidden          bool             `json:"hidden"`
	Frozen          bool             `json:"frozen"`
}
