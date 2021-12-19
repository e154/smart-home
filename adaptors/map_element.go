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
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// IMapElement ...
type IMapElement interface {
	Add(ver *m.MapElement) (id int64, err error)
	GetById(mapId int64) (ver *m.MapElement, err error)
	GetByName(name string) (ver *m.MapElement, err error)
	Update(ver *m.MapElement) (err error)
	Delete(mapId int64) (err error)
	Sort(ver *m.MapElement) (err error)
	List(limit, offset int64, orderBy, sort string) (list []*m.MapElement, total int64, err error)
	GetActiveElements(sortBy, order string, limit, offset int) (result []*m.MapElement, total int64, err error)
	fromDb(dbVer *db.MapElement) (ver *m.MapElement)
	toDb(ver *m.MapElement) (dbVer *db.MapElement)
}

// Entity ...
type MapElement struct {
	IMapElement
	table *db.MapElements
	db    *gorm.DB
}

// GetMapElementAdaptor ...
func GetMapElementAdaptor(d *gorm.DB) IMapElement {
	return &MapElement{
		table: &db.MapElements{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapElement) Add(ver *m.MapElement) (id int64, err error) {

	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		err = errors.Wrap(common.ErrTransactionError, err.Error())
		return
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(common.ErrTransactionError, err.Error())
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var protId int64
	switch {
	case ver.Prototype.MapText != nil:
		textAdaptor := GetMapTextAdaptor(tx)
		protId, err = textAdaptor.Add(ver.Prototype.MapText)
		ver.PrototypeType = common.MapElementPrototypeText
	case ver.Prototype.MapImage != nil:
		imageAdaptor := GetMapImageAdaptor(tx)
		protId, err = imageAdaptor.Add(ver.Prototype.MapImage)
		ver.PrototypeType = common.MapElementPrototypeImage
	case ver.Prototype.Entity != nil:
		ver.PrototypeId = ver.Prototype.Entity.Id
		ver.PrototypeType = common.MapElementPrototypeEntity
	default:

	}

	if protId != 0 {
		ver.PrototypeId = fmt.Sprintf("%d", protId)
	}

	if err != nil {
		return
	}

	dbVer := n.toDb(ver)
	table := db.MapElements{Db: tx}
	if id, err = table.Add(dbVer); err != nil {
		return
	}

	return
}

// GetById ...
func (n *MapElement) GetById(mapId int64) (ver *m.MapElement, err error) {

	var dbVer *db.MapElement
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// GetByName ...
func (n *MapElement) GetByName(name string) (ver *m.MapElement, err error) {

	var dbVer *db.MapElement
	if dbVer, err = n.table.GetByName(name); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

// Update ...
func (n *MapElement) Update(ver *m.MapElement) (err error) {

	var oldVer *m.MapElement
	if oldVer, err = n.GetById(ver.Id); err != nil {
		return
	}

	if oldVer.PrototypeId == "" {
		oldVer.PrototypeType = ""
	}

	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		err = errors.Wrap(common.ErrTransactionError, err.Error())
		return
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(common.ErrTransactionError, err.Error())
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	var deleted bool
	// delete old prototype
	if oldVer.PrototypeType != ver.PrototypeType {
		deleted = true
		switch oldVer.PrototypeType {
		case common.MapElementPrototypeText:
			textAdaptor := GetMapTextAdaptor(tx)
			if id, ok := oldVer.PrototypeId.(int64); ok {
				err = textAdaptor.Delete(id)
			}
		case common.MapElementPrototypeImage:
			imageAdaptor := GetMapImageAdaptor(tx)
			if id, ok := oldVer.PrototypeId.(int64); ok {
				err = imageAdaptor.Delete(id)
			}
		case common.MapElementPrototypeEntity:

		case common.MapElementPrototypeEmpty:
			log.Warn("empty prototype")
		default:
			log.Warnf("unknown prototype: '%v'", oldVer.PrototypeType)
		}

		if err != nil {
			return
		}
	}

	if ver.PrototypeId == "" {
		err = fmt.Errorf("prototype_id is zero")
		return
	}

	// add new prototype
	switch ver.PrototypeType {
	case common.MapElementPrototypeText:
		textAdaptor := GetMapTextAdaptor(tx)
		if deleted {
			// add new
			ver.PrototypeId, err = textAdaptor.Add(ver.Prototype.MapText)
		} else {
			// update
			if id, ok := ver.PrototypeId.(int64); ok {
				ver.Prototype.MapText.Id = id
			}
			err = textAdaptor.Update(ver.Prototype.MapText)
		}
	case common.MapElementPrototypeImage:
		imageAdaptor := GetMapImageAdaptor(tx)
		if deleted {
			// add new
			ver.PrototypeId, err = imageAdaptor.Add(ver.Prototype.MapImage)
		} else {
			if id, ok := ver.PrototypeId.(int64); ok {
				mapImage := &m.MapImage{
					Id:      id,
					ImageId: ver.Prototype.MapImage.ImageId,
					Style:   "", //	TODO add style to image
				}
				err = imageAdaptor.Update(mapImage)
			}
		}
	case common.MapElementPrototypeEntity:
		ver.PrototypeId = ver.Prototype.Entity.Id
		ver.PrototypeType = common.MapElementPrototypeEntity
	default:
		err = fmt.Errorf("unknown prototype: %v", ver.PrototypeType)
		log.Warnf(err.Error())
	}

	if err != nil {
		return
	}

	dbVer := n.toDb(ver)
	table := db.MapElements{Db: tx}
	if err = table.Update(dbVer); err != nil {
		return
	}

	return
}

// Delete ...
func (n *MapElement) Delete(mapId int64) (err error) {

	var ver *m.MapElement
	if ver, err = n.GetById(mapId); err != nil {
		return
	}

	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		err = errors.Wrap(common.ErrTransactionError, err.Error())
		return
	}
	defer func() {
		if err != nil {
			err = errors.Wrap(common.ErrTransactionError, err.Error())
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	if ver.PrototypeId != "" {
		switch ver.PrototypeType {
		case common.MapElementPrototypeText:
			textAdaptor := GetMapTextAdaptor(tx)
			if id, ok := ver.PrototypeId.(int64); ok {
				err = textAdaptor.Delete(id)
			}
		case common.MapElementPrototypeImage:
			imageAdaptor := GetMapImageAdaptor(tx)
			if id, ok := ver.PrototypeId.(int64); ok {
				err = imageAdaptor.Delete(id)
			}
		case common.MapElementPrototypeEntity:
			entityAdaptor := GetEntityAdaptor(tx)
			if id, ok := ver.PrototypeId.(common.EntityId); ok {
				err = entityAdaptor.Delete(id)
			}
		default:
			err = fmt.Errorf("unknown prototype: %v", ver.PrototypeType)
			log.Warnf(err.Error())
		}
	}

	if err != nil {
		return
	}

	table := &db.MapElements{Db: tx}
	if err = table.Delete(mapId); err != nil {
		return
	}

	return
}

// Sort ...
func (n *MapElement) Sort(ver *m.MapElement) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

// List ...
func (n *MapElement) List(limit, offset int64, orderBy, sort string) (list []*m.MapElement, total int64, err error) {
	var dbList []*db.MapElement
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.MapElement, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		list = append(list, ver)
	}

	return
}

// GetActiveElements ...
func (n *MapElement) GetActiveElements(sortBy, order string, limit, offset int) (result []*m.MapElement, total int64, err error) {

	var dbList []*db.MapElement
	if dbList, total, err = n.table.GetActiveElements(int64(limit), int64(offset), order, sortBy); err != nil {
		return
	}

	result = make([]*m.MapElement, 0)

	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		result = append(result, ver)
	}

	return
}

func (n *MapElement) fromDb(dbVer *db.MapElement) (ver *m.MapElement) {
	ver = &m.MapElement{
		Id:            dbVer.Id,
		Name:          dbVer.Name,
		Description:   dbVer.Description,
		PrototypeId:   dbVer.PrototypeId,
		PrototypeType: dbVer.PrototypeType,
		LayerId:       dbVer.MapLayerId,
		MapId:         dbVer.MapId,
		Weight:        dbVer.Weight,
		Status:        dbVer.Status,
		CreatedAt:     dbVer.CreatedAt,
		UpdatedAt:     dbVer.UpdatedAt,
	}

	// GraphSettings
	graphSettings, _ := dbVer.GraphSettings.MarshalJSON()
	json.Unmarshal(graphSettings, &ver.GraphSettings)

	// MapElementPrototype
	switch {
	case dbVer.Prototype.MapText != nil:
		mapTextAdaptor := GetMapTextAdaptor(n.db)
		ver.Prototype = m.MapElementPrototype{
			MapText: mapTextAdaptor.fromDb(dbVer.Prototype.MapText),
		}
	case dbVer.Prototype.MapImage != nil:
		mapImageAdaptor := GetMapImageAdaptor(n.db)
		ver.Prototype = m.MapElementPrototype{
			MapImage: mapImageAdaptor.fromDb(dbVer.Prototype.MapImage),
		}
	case dbVer.Prototype.Entity != nil:
		mapDeviceAdaptor := GetEntityAdaptor(n.db)
		ver.Prototype = m.MapElementPrototype{
			Entity: mapDeviceAdaptor.fromDb(dbVer.Prototype.Entity),
		}
	}

	return
}

func (n *MapElement) toDb(ver *m.MapElement) (dbVer *db.MapElement) {
	dbVer = &db.MapElement{
		Id:            ver.Id,
		Name:          ver.Name,
		Description:   ver.Description,
		PrototypeId:   ver.PrototypeId,
		PrototypeType: ver.PrototypeType,
		MapLayerId:    ver.LayerId,
		MapId:         ver.MapId,
		Weight:        ver.Weight,
		Status:        ver.Status,
	}

	graphSettings, _ := json.Marshal(ver.GraphSettings)
	dbVer.GraphSettings.UnmarshalJSON(graphSettings)

	return
}
