// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

package memory_app

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// EntityMemoryApp ...
	EntityMemoryApp = string("memory_app")

	// AttrAlloc ...
	AttrAlloc = "alloc"

	// AttrHeapAlloc ...
	AttrHeapAlloc = "heap_alloc"

	// AttrTotalAlloc ...
	AttrTotalAlloc = "total_alloc"

	// AttrSys ...
	AttrSys = "sys"

	// AttrNumGC ...
	AttrNumGC = "num_gc"

	// AttrLastGC ...
	AttrLastGC = "last_gc"

	// Name ...
	Name = "memory_app"

	// EntityType ...
	EntityType = "memory_app"

	Version = "0.0.1"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrAlloc: {
			Name: AttrAlloc,
			Type: common.AttributeInt,
		},
		AttrHeapAlloc: {
			Name: AttrHeapAlloc,
			Type: common.AttributeInt,
		},
		AttrTotalAlloc: {
			Name: AttrTotalAlloc,
			Type: common.AttributeInt,
		},
		AttrSys: {
			Name: AttrSys,
			Type: common.AttributeInt,
		},
		AttrNumGC: {
			Name: AttrNumGC,
			Type: common.AttributeInt,
		},
		AttrLastGC: {
			Name: AttrLastGC,
			Type: common.AttributeTime,
		},
	}
}

func NewMetrics() []*m.Metric {
	return []*m.Metric{
		{
			Name:        "memory_app",
			Description: "App metric",
			Options: m.MetricOptions{
				Items: []m.MetricOptionsItem{
					{
						Name:        "alloc",
						Description: "",
						Color:       "#C2C2C2",
						Translate:   "alloc",
						Label:       "%",
					},
				},
			},
		},
	}
}

