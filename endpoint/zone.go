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

// ZoneEndpoint ...
type ZoneEndpoint struct {
	*CommonEndpoint
}

// NewZoneEndpoint ...
func NewZoneEndpoint(common *CommonEndpoint) *ZoneEndpoint {
	return &ZoneEndpoint{
		CommonEndpoint: common,
	}
}

// Add ...
func (n *ZoneEndpoint) Add(zone *m.Zone) (result *m.Zone, errs []*validation.Error, err error) {

	_, errs = zone.Valid()
	if len(errs) > 0 {
		return
	}

	if zone.Id, err = n.adaptors.Zone.Add(zone); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates") {
			result, err = n.adaptors.Zone.GetByName(zone.Name)
		}
		return

	} else {
		result = &m.Zone{
			Id:   zone.Id,
			Name: zone.Name,
		}
	}

	return
}

// Delete ...
func (n *ZoneEndpoint) Delete(zoneName string) (err error) {

	if _, err = n.adaptors.Zone.GetByName(zoneName); err != nil {
		return
	}

	err = n.adaptors.Zone.Delete(zoneName)

	return
}

// Search ...
func (n *ZoneEndpoint) Search(query string, limit, offset int) (result []*m.Zone, total int64, err error) {

	result, total, err = n.adaptors.Zone.Search(query, limit, offset)

	return
}
