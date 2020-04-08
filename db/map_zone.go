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

package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// MapZones ...
type MapZones struct {
	Db *gorm.DB
}

// MapZone ...
type MapZone struct {
	Id   int64 `gorm:"primary_key"`
	Name string
}

// TableName ...
func (d *MapZone) TableName() string {
	return "map_zones"
}

// Add ...
func (n MapZones) Add(zone *MapZone) (id int64, err error) {
	if err = n.Db.Create(&zone).Error; err != nil {
		return
	}
	id = zone.Id
	return
}

// GetByName ...
func (n MapZones) GetByName(zoneName string) (zone *MapZone, err error) {

	zone = &MapZone{}
	err = n.Db.Model(zone).
		Where("name = ?", zoneName).
		First(&zone).
		Error

	return
}

// Search ...
func (n *MapZones) Search(query string, limit, offset int) (list []*MapZone, total int64, err error) {

	q := n.Db.Model(&MapZone{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*MapZone, 0)
	err = q.Find(&list).Error

	return
}

// Delete ...
func (n MapZones) Delete(name string) (err error) {
	if name == "" {
		err = fmt.Errorf("zero name")
		return
	}

	err = n.Db.Delete(&MapZone{}, "name = ?", name).Error
	return
}

// Clean ...
func (n MapZones) Clean() (err error) {

	err = n.Db.Exec(`delete 
from map_zones
where id not in (
    select DISTINCT me.zone_id
    from map_elements me
    where me.zone_id notnull
    )
`).Error

	return
}
