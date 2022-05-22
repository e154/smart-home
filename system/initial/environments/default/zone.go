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

package _default

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/location"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/plugins/zone"
	. "github.com/e154/smart-home/system/initial/assertions"
)

// ZoneManager ...
type ZoneManager struct {
	adaptors *adaptors.Adaptors
}

// NewZoneManager ...
func NewZoneManager(adaptors *adaptors.Adaptors) *ZoneManager {
	return &ZoneManager{
		adaptors: adaptors,
	}
}

func (n ZoneManager) addZone(name, desc string) {

	loc, err := location.GetRegionInfo()
	So(err, ShouldBeNil)

	var latitude = loc.Latitude
	var longitude = loc.Longitude

	if latitude == 0 || longitude == 0 {
		loc, err := location.GeoLocationFromIP("")
		So(err, ShouldBeNil)
		latitude = loc.Lat
		longitude = loc.Lon
	}

	attributes := m.Attributes{
		AttrLat: {
			Name:  AttrLat,
			Type:  common.AttributeFloat,
			Value: latitude,
		},
		AttrLon: {
			Name:  AttrLon,
			Type:  common.AttributeFloat,
			Value: longitude,
		},
		AttrElevation: {
			Name:  AttrElevation,
			Type:  common.AttributeFloat,
			Value: 150,
		},
		AttrTimezone: {
			Name:  AttrTimezone,
			Type:  common.AttributeInt,
			Value: 7,
		},
	}

	err = n.adaptors.Entity.Add(&m.Entity{
		Id:          common.EntityId("zone." + name),
		Description: desc,
		PluginName:  "zone",
		Attributes:  attributes,
		AutoLoad:    true,
	})
	So(err, ShouldBeNil)

	_, err = n.adaptors.EntityStorage.Add(&m.EntityStorage{
		EntityId:   common.EntityId("zone." + name),
		Attributes: attributes.Serialize(),
	})
	So(err, ShouldBeNil)

}

// Create ...
func (n ZoneManager) Create() {
	n.addZone("home", "base geo position")
}
