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

package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/e154/smart-home/pkg/apperr"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"gorm.io/gorm"
)

// Images ...
type Images struct {
	*Common
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
	CreatedAt time.Time `gorm:"<-:create"`
}

// TableName ...
func (m *Image) TableName() string {
	return "images"
}

// Add ...
func (n Images) Add(ctx context.Context, v *Image) (id int64, err error) {
	if err = n.DB(ctx).Create(&v).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "images_pkey") {
					err = fmt.Errorf("%s: %w", fmt.Sprintf("image name \"%d\" not unique", v.Id), apperr.ErrImageAdd)
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageAdd)
		return
	}
	id = v.Id
	return
}

// GetById ...
func (n Images) GetById(ctx context.Context, id int64) (v *Image, err error) {

	v = &Image{Id: id}
	if err = n.DB(ctx).First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("%s: %w", fmt.Sprintf("id \"%d\"", id), apperr.ErrImageNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageGet)
	}
	return
}

// GetByImageName ...
func (n Images) GetByImageName(ctx context.Context, imageName string) (v *Image, err error) {
	v = &Image{}
	if err = n.DB(ctx).Model(v).Where("image = ?", imageName).First(&v).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = fmt.Errorf("%s: %w", fmt.Sprintf("name \"%s\"", imageName), apperr.ErrImageNotFound)
			return
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageGet)
		return
	}
	return
}

// Update ...
func (n Images) Update(ctx context.Context, m *Image) (err error) {
	err = n.DB(ctx).Model(&Image{Id: m.Id}).Updates(map[string]interface{}{
		"title": m.Title,
		"Name":  m.Name,
	}).Error
	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageUpdate)
		return
	}
	return
}

// Delete ...
func (n Images) Delete(ctx context.Context, mapId int64) (err error) {
	if err = n.DB(ctx).Delete(&Image{Id: mapId}).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageDelete)
		return
	}
	return
}

// List ...
func (n *Images) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*Image, total int64, err error) {

	if err = n.DB(ctx).Model(Image{}).Count(&total).Error; err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageList)
		return
	}

	list = make([]*Image, 0)
	err = n.DB(ctx).
		Limit(limit).
		Offset(offset).
		Order(fmt.Sprintf("%s %s", sort, orderBy)).
		Find(&list).
		Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageList)
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
func (n *Images) GetFilterList(ctx context.Context) (images []*ImageFilterList, err error) {

	image := &Image{}
	var rows *sql.Rows
	rows, err = n.DB(ctx).Raw(`
SELECT
	to_char(created_at,'YYYY-mm-dd') as date, COUNT( created_at) as count
FROM ` + image.TableName() + `
GROUP BY date
ORDER BY date`).Rows()

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageList)
		return
	}

	for rows.Next() {
		item := &ImageFilterList{}
		_ = rows.Scan(&item.Date, &item.Count)
		images = append(images, item)
	}

	return
}

// GetAllByDate ...
func (n *Images) GetAllByDate(ctx context.Context, filter string) (images []*Image, err error) {

	//fmt.Println("filter", filter)

	images = make([]*Image, 0)
	image := &Image{}
	err = n.DB(ctx).Raw(`
SELECT *
FROM `+image.TableName()+`
WHERE to_char(created_at,'YYYY-mm-dd') = ?
ORDER BY created_at`, filter).
		Find(&images).
		Error

	if err != nil {
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageList)
	}
	return
}

// AddMultiple ...
func (n *Images) AddMultiple(ctx context.Context, images []*Image) (err error) {
	if err = n.DB(ctx).Create(&images).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "images_pkey") {
					err = fmt.Errorf("%s: %w", "multiple insert", apperr.ErrImageAdd)
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = fmt.Errorf("%s: %w", err.Error(), apperr.ErrImageAdd)
	}
	return
}
