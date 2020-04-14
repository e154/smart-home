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

package adaptors

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

// MapElement ...
type MapElement struct {
	table *db.MapElements
	db    *gorm.DB
}

// GetMapElementAdaptor ...
func GetMapElementAdaptor(d *gorm.DB) *MapElement {
	return &MapElement{
		table: &db.MapElements{Db: d},
		db:    d,
	}
}

// Add ...
func (n *MapElement) Add(ver *m.MapElement) (id int64, err error) {

	transaction := true
	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		tx = n.db
		transaction = false
	}

	defer func() {
		if err != nil && transaction {
			tx.Rollback()
		}
	}()

	switch {
	case ver.Prototype.MapText != nil:
		textAdaptor := GetMapTextAdaptor(tx)
		ver.PrototypeId, err = textAdaptor.Add(ver.Prototype.MapText)
		ver.PrototypeType = common.PrototypeTypeText
	case ver.Prototype.MapImage != nil:
		imageAdaptor := GetMapImageAdaptor(tx)
		ver.PrototypeId, err = imageAdaptor.Add(ver.Prototype.MapImage)
		ver.PrototypeType = common.PrototypeTypeImage
	case ver.Prototype.MapDevice != nil:
		deviceAdaptor := GetMapDeviceAdaptor(tx)
		if ver.PrototypeId, err = deviceAdaptor.Add(ver.Prototype.MapDevice); err != nil {
			return
		}

		ver.PrototypeType = common.PrototypeTypeDevice
		//actions
		deviceAction := GetMapDeviceActionAdaptor(tx)
		//err = deviceAction.AddMultiple(t.Actions)
		for _, action := range ver.Prototype.MapDevice.Actions {
			action.MapDeviceId = ver.PrototypeId
			if action.Id, err = deviceAction.Add(action); err != nil {
				log.Error(err.Error())
				return
			}
		}

		//states
		stateAdaptor := GetMapDeviceStateAdaptor(tx)
		//err = stateAdaptor.AddMultiple(t.States)
		for _, state := range ver.Prototype.MapDevice.States {
			state.MapDeviceId = ver.PrototypeId
			if state.Id, err = stateAdaptor.Add(state); err != nil {
				log.Error(err.Error())
				return
			}
		}
	default:

	}

	if err != nil {
		return
	}

	dbVer := n.toDb(ver)
	table := db.MapElements{Db: tx}
	if id, err = table.Add(dbVer); err != nil {
		return
	}

	if transaction {
		err = tx.Commit().Error
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

	fmt.Println("update")
	var oldVer *m.MapElement
	if oldVer, err = n.GetById(ver.Id); err != nil {
		return
	}

	if oldVer.PrototypeId == 0 {
		oldVer.PrototypeType = ""
	}

	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var deleted bool
	// delete old prototype
	if oldVer.PrototypeType != ver.PrototypeType {
		deleted = true
		switch oldVer.PrototypeType {
		case common.PrototypeTypeText:
			textAdaptor := GetMapTextAdaptor(tx)
			err = textAdaptor.Delete(oldVer.PrototypeId)
		case common.PrototypeTypeImage:
			imageAdaptor := GetMapImageAdaptor(tx)
			err = imageAdaptor.Delete(oldVer.PrototypeId)
		case common.PrototypeTypeDevice:
			deviceAdaptor := GetMapDeviceAdaptor(tx)
			err = deviceAdaptor.Delete(oldVer.PrototypeId)
		case common.PrototypeTypeEmpty:
			log.Warn("empty prototype")
		default:
			log.Warnf("unknown prototype: '%v'", oldVer.PrototypeType)
		}

		if err != nil {
			return
		}
	}

	if ver.PrototypeId == 0 {
		err = fmt.Errorf("prototype_id is zero")
		return
	}

	// add new prototype
	switch ver.PrototypeType {
	case common.PrototypeTypeText:
		textAdaptor := GetMapTextAdaptor(tx)
		if deleted {
			// add new
			ver.PrototypeId, err = textAdaptor.Add(ver.Prototype.MapText)
		} else {
			// update
			ver.Prototype.MapText.Id = ver.PrototypeId
			err = textAdaptor.Update(ver.Prototype.MapText)
		}
	case common.PrototypeTypeImage:
		imageAdaptor := GetMapImageAdaptor(tx)
		if deleted {
			// add new
			ver.PrototypeId, err = imageAdaptor.Add(ver.Prototype.MapImage)
		} else {
			mapImage := &m.MapImage{
				Id:      ver.PrototypeId,
				ImageId: ver.Prototype.MapImage.ImageId,
				Style:   "", //	TODO add style to image
			}
			err = imageAdaptor.Update(mapImage)
		}
	case common.PrototypeTypeDevice:
		deviceAdaptor := GetMapDeviceAdaptor(tx)
		if deleted {
			// add new
			ver.PrototypeId, err = deviceAdaptor.Add(ver.Prototype.MapDevice)
		} else {
			ver.PrototypeId = oldVer.PrototypeId
			ver.Prototype.MapDevice.Id = ver.PrototypeId
			if err = deviceAdaptor.Update(ver.Prototype.MapDevice); err != nil {
				log.Error(err.Error())
				return
			}

			mapDeviceActionAdaptor := GetMapDeviceActionAdaptor(tx)
			if err = mapDeviceActionAdaptor.DeleteByDeviceId(oldVer.PrototypeId); err != nil {
				return
			}

			mapDeviceStateAdaptor := GetMapDeviceStateAdaptor(tx)
			if err = mapDeviceStateAdaptor.DeleteByDeviceId(oldVer.PrototypeId); err != nil {
				return
			}
		}

		if ver.Prototype.MapDevice != nil {

			//actions
			for _, action := range ver.Prototype.MapDevice.Actions {
				action.MapDeviceId = ver.PrototypeId
			}
			deviceAction := GetMapDeviceActionAdaptor(tx)
			if err = deviceAction.AddMultiple(ver.Prototype.MapDevice.Actions); err != nil {
				log.Error(err.Error())
				return
			}

			//states
			for _, state := range ver.Prototype.MapDevice.States {
				state.MapDeviceId = ver.PrototypeId
			}
			stateAdaptor := GetMapDeviceStateAdaptor(tx)
			if err = stateAdaptor.AddMultiple(ver.Prototype.MapDevice.States); err != nil {
				log.Errorf(err.Error())
				return
			}

			//metrics
			metricAdaptor := GetMetricAdaptor(tx, nil)

			if oldVer.PrototypeType == common.PrototypeTypeDevice {
				for _, oldMetric := range oldVer.Prototype.Metrics {
					var exist bool
					for _, metric := range ver.Prototype.MapDevice.Metrics {
						if metric.Id == oldMetric.Id {
							exist = true
						}
					}
					if !exist {
						if err = metricAdaptor.Delete(oldMetric.Id); err != nil {
							return
						}
					}
				}
			}

			for _, metric := range ver.Prototype.MapDevice.Metrics {
				var exist bool
				if oldVer.PrototypeType == common.PrototypeTypeDevice {
					for _, oldMetric := range oldVer.Prototype.Metrics {
						if metric.Id == oldMetric.Id {
							exist = true
						}
					}
				}
				if !exist {
					metric.MapDeviceId = ver.PrototypeId
					if _, err = metricAdaptor.Add(metric); err != nil {
						return
					}
				} else {
					if err = metricAdaptor.Update(metric); err != nil {
						return
					}
				}
			}
		} else {
			// delete metrics
			metricAdaptor := GetMetricAdaptor(tx, nil)
			if err = metricAdaptor.DeleteByDeviceId(oldVer.PrototypeId); err != nil {
				return
			}
		}

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

	err = tx.Commit().Error

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
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if ver.PrototypeId != 0 {
		switch ver.PrototypeType {
		case common.PrototypeTypeText:
			textAdaptor := GetMapTextAdaptor(tx)
			err = textAdaptor.Delete(ver.PrototypeId)
		case common.PrototypeTypeImage:
			imageAdaptor := GetMapImageAdaptor(tx)
			err = imageAdaptor.Delete(ver.PrototypeId)
		case common.PrototypeTypeDevice:
			deviceAdaptor := GetMapDeviceAdaptor(tx)
			err = deviceAdaptor.Delete(ver.PrototypeId)
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

	err = tx.Commit().Error

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

	// Zone tag
	if dbVer.Zone != nil {
		zoneAdaptor := GetMapZoneAdaptor(n.db)
		ver.Zone = zoneAdaptor.fromDb(dbVer.Zone)
	}

	// GraphSettings
	graphSettings, _ := dbVer.GraphSettings.MarshalJSON()
	json.Unmarshal(graphSettings, &ver.GraphSettings)

	// Prototype
	switch {
	case dbVer.Prototype.MapText != nil:
		mapTextAdaptor := GetMapTextAdaptor(n.db)
		ver.Prototype = m.Prototype{
			MapText: mapTextAdaptor.fromDb(dbVer.Prototype.MapText),
		}
	case dbVer.Prototype.MapImage != nil:
		mapImageAdaptor := GetMapImageAdaptor(n.db)
		ver.Prototype = m.Prototype{
			MapImage: mapImageAdaptor.fromDb(dbVer.Prototype.MapImage),
		}
	case dbVer.Prototype.MapDevice != nil:
		mapDeviceAdaptor := GetMapDeviceAdaptor(n.db)
		ver.Prototype = m.Prototype{
			MapDevice: mapDeviceAdaptor.fromDb(dbVer.Prototype.MapDevice),
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

	if ver.Zone != nil && ver.Zone.Id != 0 {
		dbVer.ZoneId = &ver.Zone.Id
	} else {
		dbVer.ZoneId = nil
	}

	graphSettings, _ := json.Marshal(ver.GraphSettings)
	dbVer.GraphSettings.UnmarshalJSON(graphSettings)

	return
}
