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

package cpuspeed

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// EntityCpuspeed ...
	EntityCpuspeed = string("cpuspeed")

	//icon = "microchip"

	// AttrCpuCores ...
	AttrCpuCores = "cores"
	// AttrCpuMhz ...
	AttrCpuMhz = "mhz"
	// AttrCpuAll ...
	AttrCpuAll = "all"
	// AttrLoadMin ...
	AttrLoadMin = "load_min"
	// AttrLoadMax ...
	AttrLoadMax = "load_max"

	// Name ...
	Name = "cpuspeed"

	// EntityType ...
	EntityType = "cpuspeed"

	Version = "0.0.1"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return m.Attributes{
		AttrCpuCores: {
			Name: AttrCpuCores,
			Type: common.AttributeInt,
		},
		AttrCpuMhz: {
			Name: AttrCpuMhz,
			Type: common.AttributeFloat,
		},
		AttrCpuAll: {
			Name: AttrCpuAll,
			Type: common.AttributeFloat,
		},
		AttrLoadMin: {
			Name: AttrLoadMin,
			Type: common.AttributeFloat,
		},
		AttrLoadMax: {
			Name: AttrLoadMax,
			Type: common.AttributeFloat,
		},
	}
}

func NewMetrics() []*m.Metric {
	return []*m.Metric{
		{
			Name:        "cpuspeed",
			Description: "Cpu metric",
			Options: m.MetricOptions{
				Items: []m.MetricOptionsItem{
					{
						Name:        "all",
						Description: "",
						Color:       "#C2C2C2",
						Translate:   "all",
						Label:       "GHz",
					},
				},
			},
		},
	}
}
