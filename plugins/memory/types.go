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

package memory

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// EntityMemory ...
	EntityMemory = string("memory")

	// AttrTotal ...
	AttrTotal = "total"
	// AttrFree ...
	AttrFree = "free"
	// AttrUsedPercent ...
	AttrUsedPercent = "used_percent"

	// Name ...
	Name = "memory"

	// EntityType ...
	EntityType = "memory"

	Version = "0.0.1"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrTotal: {
			Name: AttrTotal,
			Type: common.AttributeInt,
		},
		AttrFree: {
			Name: AttrFree,
			Type: common.AttributeInt,
		},
		AttrUsedPercent: {
			Name: AttrUsedPercent,
			Type: common.AttributeFloat,
		},
	}
}

func NewMetrics() []*m.Metric {
	return []*m.Metric{
		{
			Name:        "memory",
			Description: "RAM metric",
			Options: m.MetricOptions{
				Items: []m.MetricOptionsItem{
					{
						Name:        "used_percent",
						Description: "",
						Color:       "#C2C2C2",
						Translate:   "used_percent",
						Label:       "%",
					},
				},
			},
		},
	}
}
