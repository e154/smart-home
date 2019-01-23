package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"fmt"
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

	switch t := ver.Prototype.(type) {
	case *m.MapText:
		textAdaptor := GetMapTextAdaptor(n.db)
		ver.PrototypeId, err = textAdaptor.Add(t)
		ver.PrototypeType = "text"
	case *m.MapImage:
		imageAdaptor := GetMapImageAdaptor(n.db)
		ver.PrototypeId, err = imageAdaptor.Add(t)
		ver.PrototypeType = "image"
	case *m.MapDevice:
		deviceAdaptor := GetMapDeviceAdaptor(n.db)
		if ver.PrototypeId, err = deviceAdaptor.Add(t); err != nil {
			return
		}
		ver.PrototypeType = "device"
		//actions
		deviceAction := GetMapDeviceActionAdaptor(n.db)
		deviceAction.AddMultiple(t.Actions)
		//states
		stateAdaptor := GetMapDeviceStateAdaptor(n.db)
		stateAdaptor.AddMultiple(t.States)

	default:
		err = fmt.Errorf("unknown prototype: %v", ver.PrototypeType)
		log.Warningf(err.Error())
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

	// delete old prototype
	switch oldVer.PrototypeType {
	case "text":
		textAdaptor := GetMapTextAdaptor(n.db)
		err = textAdaptor.Delete(ver.PrototypeId)
	case "image":
		imageAdaptor := GetMapImageAdaptor(n.db)
		err = imageAdaptor.Delete(ver.PrototypeId)
	case "device":
		deviceAdaptor := GetMapDeviceAdaptor(n.db)
		err = deviceAdaptor.Delete(ver.PrototypeId)
	default:
		err = fmt.Errorf("unknown prototype: %v", ver.PrototypeType)
		log.Warningf(err.Error())

	}

	if err != nil {
		return
	}

	// add new prototype
	switch t := ver.Prototype.(type) {
	case *m.MapText:
		textAdaptor := GetMapTextAdaptor(n.db)
		ver.PrototypeId, err = textAdaptor.Add(t)
		ver.PrototypeType = "text"
	case *m.MapImage:
		imageAdaptor := GetMapImageAdaptor(n.db)
		ver.PrototypeId, err = imageAdaptor.Add(t)
		ver.PrototypeType = "image"
	case *m.MapDevice:
		deviceAdaptor := GetMapDeviceAdaptor(n.db)
		if ver.PrototypeId, err = deviceAdaptor.Add(t); err != nil {
			return
		}
		ver.PrototypeType = "device"
		//actions
		deviceAction := GetMapDeviceActionAdaptor(n.db)
		deviceAction.AddMultiple(t.Actions)
		//states
		stateAdaptor := GetMapDeviceStateAdaptor(n.db)
		stateAdaptor.AddMultiple(t.States)
	default:
		log.Warningf("unknown prototype: %v", t)
	}

	if err != nil {
		return
	}

	dbVer := n.toDb(ver)
	err = n.table.Update(dbVer)
	return
}

func (n *MapElement) Delete(mapId int64) (err error) {

	var ver *m.MapElement
	if ver, err = n.GetById(mapId); err != nil {
		return
	}

	switch ver.PrototypeType {
	case "text":
		textAdaptor := GetMapTextAdaptor(n.db)
		err = textAdaptor.Delete(ver.PrototypeId)
	case "image":
		imageAdaptor := GetMapImageAdaptor(n.db)
		err = imageAdaptor.Delete(ver.PrototypeId)
	case "device":
		deviceAdaptor := GetMapDeviceAdaptor(n.db)
		err = deviceAdaptor.Delete(ver.PrototypeId)
	default:
		err = fmt.Errorf("unknown prototype: %v", ver.PrototypeType)
		log.Warningf(err.Error())
	}

	if err != nil {
		return
	}

	err = n.table.Delete(mapId)
	return
}

func (n *MapElement) Sort(ver *m.MapElement) (err error) {
	dbVer := n.toDb(ver)
	err = n.table.Sort(dbVer)
	return
}

func (n *MapElement) fromDb(dbVer *db.MapElement) (ver *m.MapElement) {
	ver = &m.MapElement{
		Id:            dbVer.Id,
		Name:          dbVer.Name,
		Description:   dbVer.Description,
		PrototypeId:   dbVer.PrototypeId,
		PrototypeType: dbVer.PrototypeType,
		LayerId:       dbVer.LayerId,
		MapId:         dbVer.MapId,
		Weight:        dbVer.Weight,
		GraphSettings: dbVer.GraphSettings,
		Status:        dbVer.Status,
		CreatedAt:     dbVer.CreatedAt,
		UpdatedAt:     dbVer.UpdatedAt,
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
		LayerId:       ver.LayerId,
		MapId:         ver.MapId,
		Weight:        ver.Weight,
		GraphSettings: ver.GraphSettings,
		Status:        ver.Status,
	}
	return
}
