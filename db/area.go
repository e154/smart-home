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
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/e154/smart-home/common/apperr"
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
	Polygon     *Polygon
	Payload     json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName ...
func (d *Area) TableName() string {
	return "areas"
}

// Add ...
func (n *Areas) Add(ctx context.Context, area *Area) (id int64, err error) {
	if err = n.Db.WithContext(ctx).Create(&area).Error; err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				if strings.Contains(pgErr.Message, "name_at_areas_unq") {
					err = errors.Wrap(apperr.ErrAreaAdd, fmt.Sprintf("area name \"%s\" not unique", area.Name))
					return
				}
			default:
				fmt.Printf("unknown code \"%s\"\n", pgErr.Code)
			}
		}
		err = errors.Wrap(apperr.ErrAreaAdd, err.Error())
		return
	}
	id = area.Id
	return
}

// GetByName ...
func (n *Areas) GetByName(ctx context.Context, name string) (area *Area, err error) {

	area = &Area{}
	err = n.Db.WithContext(ctx).Model(area).
		Where("name = ?", name).
		First(&area).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrAreaGet, err.Error())
	}
	return
}

// Search ...
func (n *Areas) Search(ctx context.Context, query string, limit, offset int) (list []*Area, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Area{}).
		Where("name LIKE ?", "%"+query+"%")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrAreaGet, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	list = make([]*Area, 0)
	if err = q.Find(&list).Error; err != nil {
		err = errors.Wrap(apperr.ErrAreaGet, err.Error())
	}

	return
}

// DeleteByName ...
func (n *Areas) DeleteByName(ctx context.Context, name string) (err error) {
	if name == "" {
		err = errors.Wrap(apperr.ErrAreaDelete, "zero name")
		return
	}

	if err = n.Db.WithContext(ctx).Delete(&Area{}, "name = ?", name).Error; err != nil {
		err = errors.Wrap(apperr.ErrAreaDelete, "zero name")
	}
	return
}

// Clean ...
func (n *Areas) Clean(ctx context.Context) (err error) {

	err = n.Db.WithContext(ctx).Exec(`delete 
from areas
where id not in (
    select DISTINCT me.area_id
    from entities me
    where me.area_id notnull
    )
`).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrAreaClean, "zero name")
	}

	return
}

// Update ...
func (n *Areas) Update(ctx context.Context, m *Area) (err error) {
	err = n.Db.WithContext(ctx).Model(&Area{Id: m.Id}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"payload":     m.Payload,
		"polygon":     m.Polygon,
	}).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrAreaUpdate, err.Error())
	}
	return
}

// List ...
func (n *Areas) List(ctx context.Context, limit, offset int, orderBy, sort string) (list []*Area, total int64, err error) {

	if err = n.Db.WithContext(ctx).Model(Area{}).Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrAreaList, err.Error())
		return
	}

	list = make([]*Area, 0)
	q := n.Db.WithContext(ctx).Model(&Area{}).
		Limit(limit).
		Offset(offset)

	if sort != "" && orderBy != "" {
		q = q.
			Order(fmt.Sprintf("%s %s", sort, orderBy))
	}

	err = q.
		Find(&list).
		Error
	if err != nil {
		err = errors.Wrap(apperr.ErrAreaList, err.Error())
	}

	return
}

// GetById ...
func (n *Areas) GetById(ctx context.Context, areaId int64) (area *Area, err error) {
	area = &Area{Id: areaId}
	if err = n.Db.WithContext(ctx).First(&area).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrAreaNotFound, fmt.Sprintf("id \"%d\"", areaId))
			return
		}
		err = errors.Wrap(apperr.ErrAreaGet, err.Error())
	}
	return
}

func (a *Areas) ListByPoint(ctx context.Context, point Point, limit, offset int) (list []*Area, err error) {

	// https://postgis.net/docs/ST_Point.html
	// geometry ST_Point(float x, float y);
	// For geodetic coordinates, X is longitude and Y is latitude

	const query = `
SELECT *
		FROM areas as a
WHERE ST_Contains(a.polygon::geometry,
		ST_Transform(
			ST_GeomFromText('POINT(%f %f)', 4326), 4326
		)
	)`

	list = make([]*Area, 0)
	q := fmt.Sprintf(query, point.Lon, point.Lat)

	err = a.Db.WithContext(ctx).Raw(q).
		Limit(limit).
		Offset(offset).Scan(&list).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrAreaList, err.Error())
		return
	}

	return
}

func (a *Areas) GetDistance(ctx context.Context, point Point, areaId int64) (distance float64, err error) {

	const query = `
select st_distance(
   ST_Transform(ST_GeomFromText('POINT (%f %f)', 4326)::geometry, 4326),
   ST_Transform((select polygon from areas where id = %d)::geometry, 4326)
)`
	q := fmt.Sprintf(query, point.Lat, point.Lon, areaId)
	err = a.Db.WithContext(ctx).Raw(q).Scan(&distance).Error
	if err != nil {
		err = errors.Wrap(apperr.ErrAreaList, err.Error())
		return
	}

	return
}
