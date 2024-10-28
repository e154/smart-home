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

// DashboardTab ...
type DashboardTab struct {
	Payload     json.RawMessage             `json:"payload"`
	Cards       []*DashboardCard            `json:"cards"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
	Name        string                      `json:"name" validate:"required"`
	Icon        string                      `json:"icon"`
	Id          int64                       `json:"id"`
	ColumnWidth int                         `json:"column_width"`
	Background  *string                     `json:"background"`
	Weight      int                         `json:"weight"`
	DashboardId int64                       `json:"dashboard_id" validate:"required"`
	Dashboard   *Dashboard                  `json:"dashboard"`
	Entities    map[common.EntityId]*Entity `json:"entities"`
	Gap         bool                        `json:"gap"`
	Enabled     bool                        `json:"enabled"`
}
