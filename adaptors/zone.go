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

package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type IZone interface {
	Add(tag *m.Zone) (id int64, err error)
	GetByName(zoneName string) (ver *m.Zone, err error)
	Delete(name string) (err error)
	Search(query string, limit, offset int) (list []*m.Zone, total int64, err error)
	Clean() (err error)
	toDb(tag *m.Zone) *db.Zone
	fromDb(tag *db.Zone) *m.Zone
}

// Zone ...
type Zone struct {
	IZone
	table *db.Zones
	db    *gorm.DB
}

// GetZoneAdaptor ...
func GetZoneAdaptor(d *gorm.DB) IZone {
	return &Zone{
		table: &db.Zones{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Zone) Add(tag *m.Zone) (id int64, err error) {

	dbTag := n.toDb(tag)
	id, err = n.table.Add(dbTag)

	return
}

// GetByName ...
func (n *Zone) GetByName(zoneName string) (ver *m.Zone, err error) {

	var dbVer *db.Zone
	if dbVer, err = n.table.GetByName(zoneName); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Delete ...
func (n *Zone) Delete(name string) (err error) {

	err = n.table.Delete(name)

	return
}

// Search ...
func (n *Zone) Search(query string, limit, offset int) (list []*m.Zone, total int64, err error) {
	var dbList []*db.Zone
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Zone, 0)
	for _, dbTag := range dbList {
		node := n.fromDb(dbTag)
		list = append(list, node)
	}

	return
}

// Clean ...
func (n *Zone) Clean() (err error) {

	err = n.table.Clean()

	return
}

func (n *Zone) toDb(tag *m.Zone) *db.Zone {
	return &db.Zone{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func (n *Zone) fromDb(tag *db.Zone) *m.Zone {
	return &m.Zone{
		Id:   tag.Id,
		Name: tag.Name,
	}
}
