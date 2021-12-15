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

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/e154/smart-home/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Images ...
type Images struct {
	Db *gorm.DB
}

// Image ...
type Image struct {
	Id        int64 `gorm:"primary_key"`
	Thumb     string
	Image     string
	MimeType  string
	Title     string
	Size      int64
	Name      string
	CreatedAt time.Time
}

// TableName ...
func (m *Image) TableName() string {
	return "images"
}

// Add ...
func (n Images) Add(v *Image) (id int64, err error) {
	if err = n.Db.Create(&v).Error; err != nil {
		err = errors.Wrap(err, "add failed")
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n Images) GetById(id int64) (v *Image, err error) {
	v = &Image{Id: id}
	if err = n.Db.First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(common.ErrNotFound, fmt.Sprintf("map with id \"%d\"", id))
		}
	}
	return
}

// GetByImageName ...
func (n Images) GetByImageName(imageName string) (v *Image, err error) {
	v = &Image{}
	if err = n.Db.Model(v).Where("image = ?", imageName).First(&v).Error; err != nil {
		err = errors.Wrap(err, "getByImageName failed")
		return
	}
	return
}

// Update ...
func (n Images) Update(m *Image) (err error) {
	err = n.Db.Model(&Image{Id: m.Id}).Updates(map[string]interface{}{
		"title": m.Title,
		"Name":  m.Name,
	}).Error
	if err != nil {
		err = errors.Wrap(err, "update failed")
		return
	}
	return
}

// Delete ...
func (n Images) Delete(mapId int64) (err error) {
	if err = n.Db.Delete(&Image{Id: mapId}).Error; err != nil {
		err = errors.Wrap(err, "delete failed")
		return
	}
	return
}

// List ...
func (n *Images) List(limit, offset int64, orderBy, sort string) (list []*Image, total int64, err error) {

	if err = n.Db.Model(Image{}).Count(&total).Error; err != nil {
		err = errors.Wrap(err, "get count failed")
		return
	}

	list = make([]*Image, 0)
	err = n.Db.
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	if err != nil {
		err = errors.Wrap(err, "list failed")
		return
	}
	return
}

// ImageFilterList ...
type ImageFilterList struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetFilterList ...
func (n *Images) GetFilterList() (images []*ImageFilterList, err error) {

	image := &Image{}
	var rows *sql.Rows
	rows, err = n.Db.Raw(`
SELECT
	to_char(created_at,'YYYY-mm-dd') as date, COUNT( created_at) as count
FROM ` + image.TableName() + `
GROUP BY date
ORDER BY date`).Rows()

	if err != nil {
		err = errors.Wrap(err, "getFilterList failed")
		return
	}

	for rows.Next() {
		item := &ImageFilterList{}
		rows.Scan(&item.Date, &item.Count)
		images = append(images, item)
	}

	return
}

// GetAllByDate ...
func (n *Images) GetAllByDate(filter string) (images []*Image, err error) {

	//fmt.Println("filter", filter)

	images = make([]*Image, 0)
	image := &Image{}
	err = n.Db.Raw(`
SELECT *
FROM `+image.TableName()+`
WHERE to_char(created_at,'YYYY-mm-dd') = ?
ORDER BY created_at`, filter).
		Find(&images).
		Error

	if err != nil {
		err = errors.Wrap(err, "getAllByDate failed")
	}
	return
}
