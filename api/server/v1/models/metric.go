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

package models

import (
	"time"
)

type MetricOptionsItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Translate   string `json:"translate"`
	Label       string `json:"label"`
}

type MetricOptions struct {
	Items []MetricOptionsItem `json:"items"`
}

// swagger:model
type Metric struct {
	Id          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Options     MetricOptions         `json:"options"`
	Data        []MetricDataItemShort `json:"data"`
	Type        string                `json:"type"`
	Ranges      []string              `json:"ranges"`
	UpdatedAt   time.Time             `json:"updated_at"`
	CreatedAt   time.Time             `json:"created_at"`
}

// swagger:model
type NewMetric struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Options     MetricOptions `json:"options"`
	Type        string        `json:"type"`
}

// swagger:model
type UpdateMetric struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Options     MetricOptions `json:"options"`
	Type        string        `json:"type"`
}