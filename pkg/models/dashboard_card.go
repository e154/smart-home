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

	"github.com/e154/smart-home/pkg/common"
)

// DashboardCard ...
type DashboardCard struct {
	Payload        json.RawMessage             `json:"payload"`
	Items          []*DashboardCardItem        `json:"items"`
	CreatedAt      time.Time                   `json:"created_at"`
	UpdatedAt      time.Time                   `json:"updated_at"`
	Title          string                      `json:"title" validate:"required"`
	Id             int64                       `json:"id"`
	Height         int                         `json:"height" validate:"required"`
	Width          int                         `json:"width" validate:"required"`
	Background     *string                     `json:"background"`
	Weight         int                         `json:"weight"`
	DashboardTabId int64                       `json:"dashboard_tab_id" validate:"required"`
	DashboardTab   *DashboardTab               `json:"dashboard_tab"`
	Entities       map[common.EntityId]*Entity `json:"entities"`
	EntityId       *common.EntityId            `json:"entity_id"`
	Entity         *Entity                     `json:"entity"`
	Enabled        bool                        `json:"enabled"`
	Hidden         bool                        `json:"hidden"`
}
