// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package zone

import (
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

const (
	// Name ...
	Name = "zone"
	// EntityZone ...
	EntityZone = string(Name)
)

const (
	// AttrLat ...
	AttrLat = "lat"
	// AttrLon ...
	AttrLon = "lon"
	// AttrElevation ...
	AttrElevation = "elevation"
	// AttrTimezone ...
	AttrTimezone = "timezone"
)

// NewAttr ...
func NewAttr() m.Attributes {
	return nil
}

// NewSettings ...
func NewSettings() m.Attributes {

	attributes := m.Attributes{
		AttrLat: {
			Name: AttrLat,
			Type: common.AttributeFloat,
		},
		AttrLon: {
			Name: AttrLon,
			Type: common.AttributeFloat,
		},
		AttrElevation: {
			Name: AttrElevation,
			Type: common.AttributeFloat,
		},
		AttrTimezone: {
			Name: AttrTimezone,
			Type: common.AttributeInt,
		},
	}

	return attributes
}
