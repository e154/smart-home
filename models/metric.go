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
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	"time"
)

// MetricOptionsItem ...
type MetricOptionsItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Translate   string `json:"translate"`
	Label       string `json:"label"`
}

// MetricOptions ...
type MetricOptions struct {
	Items []MetricOptionsItem `json:"items"`
}

// Metric ...
type Metric struct {
	Id          int64             `json:"id"`
	Name        string            `json:"name" valid:"MaxSize(254);Required"`
	Description string            `json:"description" valid:"MaxSize(254)"`
	Options     MetricOptions     `json:"options"`
	Data        []MetricDataItem  `json:"data"`
	Type        common.MetricType `json:"type" valid:"Required"`
	Ranges      []string          `json:"ranges"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Valid ...
func (d *Metric) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

func (d *Metric) RangesByType() []string {

	switch d.Type {
	case common.MetricTypeLine, common.MetricTypeBar, common.MetricTypeHorizontalBar:
		d.Ranges = []string{"1h", "6h", "12h", "24h", "7d", "30d"}
	case common.MetricTypeDoughnut, common.MetricTypeRadar, common.MetricTypePie:
		d.Ranges = []string{"current"}
	default:
		d.Ranges = []string{"1h", "6h", "12h", "24h", "7d", "30d"}
	}

	return d.Ranges
}