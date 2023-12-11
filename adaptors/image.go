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

package adaptors

import (
	"context"

	"gorm.io/gorm"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IImage ...
type IImage interface {
	Add(ctx context.Context, ver *m.Image) (id int64, err error)
	GetByImageName(ctx context.Context, imageName string) (ver *m.Image, err error)
	GetById(ctx context.Context, mapId int64) (ver *m.Image, err error)
	Update(ctx context.Context, ver *m.Image) (err error)
	Delete(ctx context.Context, mapId int64) (err error)
	List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Image, total int64, err error)
	AddMultiple(ctx context.Context, items []*m.Image) (err error)
	GetAllByDate(ctx context.Context, filter string) (images []*m.Image, err error)
	GetFilterList(ctx context.Context) (filterList []*m.ImageFilterList, err error)
	fromDb(dbImage *db.Image) (image *m.Image)
	toDb(image *m.Image) (dbImage *db.Image)
}

// Image ...
type Image struct {
	IImage
	table *db.Images
	db    *gorm.DB
}

// GetImageAdaptor ...
func GetImageAdaptor(d *gorm.DB) IImage {
	return &Image{
		table: &db.Images{Db: d},
		db:    d,
	}
}

// Add ...
func (n *Image) Add(ctx context.Context, ver *m.Image) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(ctx, dbVer); err != nil {
		return
	}

	return
}

// GetByImageName ...
func (n *Image) GetByImageName(ctx context.Context, imageName string) (ver *m.Image, err error) {

	var dbVer *db.Image
	if dbVer, err = n.table.GetByImageName(ctx, imageName); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// GetById ...
func (n *Image) GetById(ctx context.Context, mapId int64) (ver *m.Image, err error) {

	var dbVer *db.Image
	if dbVer, err = n.table.GetById(ctx, mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *Image) Update(ctx context.Context, ver *m.Image) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Update(ctx, dbVer)
	return
}

// Delete ...
func (n *Image) Delete(ctx context.Context, mapId int64) (err error) {
	err = n.table.Delete(ctx, mapId)
	return
}

// List ...
func (n *Image) List(ctx context.Context, limit, offset int64, orderBy, sort string) (list []*m.Image, total int64, err error) {

	if sort == "" {
		sort = "id"
	}
	if orderBy == "" {
		orderBy = "desc"
	}

	var dbList []*db.Image
	if dbList, total, err = n.table.List(ctx, int(limit), int(offset), orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Image, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// AddMultiple ...
func (n *Image) AddMultiple(ctx context.Context, items []*m.Image) (err error) {

	insertRecords := make([]*db.Image, 0)
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = n.table.AddMultiple(ctx, insertRecords)

	return
}

// GetAllByDate ...
func (n *Image) GetAllByDate(ctx context.Context, filter string) (images []*m.Image, err error) {

	var dblist []*db.Image
	if dblist, err = n.table.GetAllByDate(ctx, filter); err != nil {
		return
	}
	for _, dbVer := range dblist {
		ver := n.fromDb(dbVer)
		images = append(images, ver)
	}

	return
}

// GetFilterList ...
func (n *Image) GetFilterList(ctx context.Context) (filterList []*m.ImageFilterList, err error) {

	var dblist []*db.ImageFilterList
	if dblist, err = n.table.GetFilterList(ctx); err != nil {
		return
	}
	for _, dbVer := range dblist {
		ver := &m.ImageFilterList{
			Date:  dbVer.Date,
			Count: dbVer.Count,
		}
		filterList = append(filterList, ver)
	}
	return
}

func (n *Image) fromDb(dbImage *db.Image) (image *m.Image) {
	image = &m.Image{
		Id:        dbImage.Id,
		Thumb:     dbImage.Thumb,
		Image:     dbImage.Image,
		MimeType:  dbImage.MimeType,
		Title:     dbImage.Title,
		Size:      dbImage.Size,
		Name:      dbImage.Name,
		CreatedAt: dbImage.CreatedAt,
	}
	if image.Image != "" {
		image.Url = common.GetLinkPath(image.Image)
	}
	return
}

func (n *Image) toDb(image *m.Image) (dbImage *db.Image) {
	dbImage = &db.Image{
		Id:       image.Id,
		Thumb:    image.Thumb,
		Image:    image.Image,
		MimeType: image.MimeType,
		Title:    image.Title,
		Size:     image.Size,
		Name:     image.Name,
	}
	return
}
