package adaptors

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
)

type MapElement struct {
	table *db.MapElements
	db    *gorm.DB
}

func GetMapElementAdaptor(d *gorm.DB) *MapElement {
	return &MapElement{
		table: &db.MapElements{Db: d},
		db:    d,
	}
}

func (n *MapElement) Add(ver *m.MapElement) (id int64, err error) {

	switch {
	case ver.Prototype.MapText != nil:
		textAdaptor := GetMapTextAdaptor(n.db)
		ver.PrototypeId, err = textAdaptor.Add(ver.Prototype.MapText)
		ver.PrototypeType = common.PrototypeTypeText
	case ver.Prototype.MapImage != nil:
		imageAdaptor := GetMapImageAdaptor(n.db)
		ver.PrototypeId, err = imageAdaptor.Add(ver.Prototype.MapImage)
		ver.PrototypeType = common.PrototypeTypeImage
	case ver.Prototype.MapDevice != nil:
		deviceAdaptor := GetMapDeviceAdaptor(n.db)
		if ver.PrototypeId, err = deviceAdaptor.Add(ver.Prototype.MapDevice); err != nil {
			return
		}

		ver.PrototypeType = common.PrototypeTypeDevice
		//actions
		deviceAction := GetMapDeviceActionAdaptor(n.db)
		//err = deviceAction.AddMultiple(t.Actions)
		for _, action := range ver.Prototype.MapDevice.Actions {
			action.MapDeviceId = ver.PrototypeId
			if action.Id, err = deviceAction.Add(action); err != nil {
				log.Error(err.Error())
			}
		}

		//states
		stateAdaptor := GetMapDeviceStateAdaptor(n.db)
		//err = stateAdaptor.AddMultiple(t.States)
		for _, state := range ver.Prototype.MapDevice.States {
			state.MapDeviceId = ver.PrototypeId
			if state.Id, err = stateAdaptor.Add(state); err != nil {
				log.Error(err.Error())
			}
		}
	default:

	}

	if err != nil {
		return
	}

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapElement) GetById(mapId int64) (ver *m.MapElement, err error) {

	var dbVer *db.MapElement
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapElement) Update(ver *m.MapElement) (err error) {

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

	// delete old prototype
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
		log.Warning("empty prototype")
	default:
		log.Warningf("unknown prototype: '%v'", oldVer.PrototypeType)
	}

	if err != nil {
		return
	}

	// add new prototype
	switch ver.PrototypeType {
	case common.PrototypeTypeText:
		textAdaptor := GetMapTextAdaptor(tx)
		ver.PrototypeId, err = textAdaptor.Add(ver.Prototype.MapText)
	case common.PrototypeTypeImage:
		imageAdaptor := GetMapImageAdaptor(tx)
		mapImage := &m.MapImage{
			ImageId: ver.Prototype.MapImage.ImageId,
			Style:   "", //	TODO add style to image
		}
		ver.PrototypeId, err = imageAdaptor.Add(mapImage)
	case common.PrototypeTypeDevice:
		deviceAdaptor := GetMapDeviceAdaptor(tx)
		if ver.PrototypeId, err = deviceAdaptor.Add(ver.Prototype.MapDevice); err != nil {
			log.Error(err.Error())
			tx.Rollback()
			return
		}

		if ver.Prototype.MapDevice != nil {

			//actions
			for _, action := range ver.Prototype.MapDevice.Actions {
				action.MapDeviceId = ver.PrototypeId
			}
			deviceAction := GetMapDeviceActionAdaptor(tx)
			if err = deviceAction.AddMultiple(ver.Prototype.MapDevice.Actions); err != nil {
				log.Error(err.Error())
				tx.Rollback()
				return
			}

			//states
			for _, state := range ver.Prototype.MapDevice.States {
				state.MapDeviceId = ver.PrototypeId
			}
			stateAdaptor := GetMapDeviceStateAdaptor(tx)
			if err = stateAdaptor.AddMultiple(ver.Prototype.MapDevice.States); err != nil {
				log.Errorf(err.Error())
				tx.Rollback()
				return
			}
		}

	default:
		err = fmt.Errorf("unknown prototype: %v", ver.PrototypeType)
		log.Warningf(err.Error())
	}

	if err != nil {
		tx.Rollback()
		return
	}

	dbVer := n.toDb(ver)
	table := db.MapElements{Db: tx}
	if err = table.Update(dbVer); err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		return
	}

	return
}

func (n *MapElement) Delete(mapId int64) (err error) {

	var ver *m.MapElement
	if ver, err = n.GetById(mapId); err != nil {
		return
	}

	tx := n.db.Begin()
	if err = tx.Error; err != nil {
		return
	}

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
			log.Warningf(err.Error())
		}
	}

	if err != nil {
		tx.Rollback()
		return
	}

	table := &db.MapElements{Db: tx}
	if err = table.Delete(mapId); err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		return
	}
	return
}

func (n *MapElement) Sort(ver *m.MapElement) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

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

	graphSettings, _ := json.Marshal(ver.GraphSettings)
	dbVer.GraphSettings.UnmarshalJSON(graphSettings)

	return
}
