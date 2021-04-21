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
	"time"
)

// Areas ...
type Areas struct {
	Db *gorm.DB
}

// Area ...
type Area struct {
	Id          int64 `gorm:"primary_key"`
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Area) TableName() string {
	return "areas"
}

// Add ...
func (n Areas) Add(area *Area) (id int64, err error) {
	if err = n.Db.Create(&area).Error; err != nil {
		return
	}
	id = area.Id
	return
}

// GetByName ...
func (n Areas) GetByName(name string) (area *Area, err error) {

	area = &Area{}
	err = n.Db.Model(area).
		Where("name = ?", name).
		First(&area).
		Error

	return
}

// Search ...
func (n *Areas) Search(query string, limit, offset int) (list []*Area, total int64, err error) {

	q := n.Db.Model(&Area{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		return
	}

	list = make([]*Area, 0)
	err = q.Find(&list).Error

	return
}

// Delete ...
func (n Areas) DeleteByName(name string) (err error) {
	if name == "" {
		err = fmt.Errorf("zero name")
		return
	}

	err = n.Db.Delete(&Area{}, "name = ?", name).Error
	return
}

// Clean ...
func (n Areas) Clean() (err error) {

	err = n.Db.Exec(`delete 
from areas
where id not in (
    select DISTINCT me.area_id
    from entities me
    where me.area_id notnull
    )
`).Error

	return
}

// Update ...
func (n Areas) Update(m *Area) (err error) {
	err = n.Db.Model(&Area{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
	}).Error
	return
}

// List ...
func (n *Areas) List(limit, offset int64, orderBy, sort string) (list []*Area, total int64, err error) {

	if err = n.Db.Model(Area{}).Count(&total).Error; err != nil {
		return
	}

	list = make([]*Area, 0)
	q := n.Db.Model(&Area{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error

	return
}

// GetById ...
func (n Areas) GetById(areaId int64) (area *Area, err error) {
	area = &Area{Id: areaId}
	err = n.Db.First(&area).Error
	return
}
