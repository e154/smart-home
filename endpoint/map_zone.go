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

package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"strings"
)

type MapZoneEndpoint struct {
	*CommonEndpoint
}

func NewMapZoneEndpoint(common *CommonEndpoint) *MapZoneEndpoint {
	return &MapZoneEndpoint{
		CommonEndpoint: common,
	}
}

func (n *MapZoneEndpoint) Add(zone *m.MapZone) (result *m.MapZone, errs []*validation.Error, err error) {

	_, errs = zone.Valid()
	if len(errs) > 0 {
		return
	}

	if zone.Id, err = n.adaptors.MapZone.Add(zone); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates") {
			result, err = n.adaptors.MapZone.GetByName(zone.Name)
		}
		return

	} else {
		result = &m.MapZone{
			Id:   zone.Id,
			Name: zone.Name,
		}
	}

	return
}

func (n *MapZoneEndpoint) Delete(zoneName string) (err error) {

	if _, err = n.adaptors.MapZone.GetByName(zoneName); err != nil {
		return
	}

	err = n.adaptors.MapZone.Delete(zoneName)

	return
}

func (n *MapZoneEndpoint) Search(query string, limit, offset int) (result []*m.MapZone, total int64, err error) {

	result, total, err = n.adaptors.MapZone.Search(query, limit, offset)

	return
}
