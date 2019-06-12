package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type MapDeviceState struct {
	table *db.MapDeviceStates
	db    *gorm.DB
}

func GetMapDeviceStateAdaptor(d *gorm.DB) *MapDeviceState {
	return &MapDeviceState{
		table: &db.MapDeviceStates{Db: d},
		db:    d,
	}
}

func (n *MapDeviceState) Add(ver *m.MapDeviceState) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapDeviceState) AddMultiple(items []*m.MapDeviceState) (err error) {

	for _, ver := range items {
		dbVer := n.toDb(ver)
		if _, err = n.table.Add(dbVer); err != nil {
			return
		}
	}

	return
}

func (n *MapDeviceState) fromDb(dbVer *db.MapDeviceState) (ver *m.MapDeviceState) {
	ver = &m.MapDeviceState{
		Id:            dbVer.Id,
		DeviceStateId: dbVer.DeviceStateId,
		MapDeviceId:   dbVer.MapDeviceId,
		ImageId:       dbVer.ImageId,
		Style:         dbVer.Style,
		CreatedAt:     dbVer.CreatedAt,
		UpdatedAt:     dbVer.UpdatedAt,
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	// state
	if dbVer.DeviceState != nil {
		stateAdaptor := GetDeviceStateAdaptor(n.db)
		ver.DeviceState = stateAdaptor.fromDb(dbVer.DeviceState)
	}

	return
}

func (n *MapDeviceState) toDb(ver *m.MapDeviceState) (dbVer *db.MapDeviceState) {
	dbVer = &db.MapDeviceState{
		Id:            ver.Id,
		DeviceStateId: ver.DeviceStateId,
		MapDeviceId:   ver.MapDeviceId,
		ImageId:       ver.ImageId,
		Style:         ver.Style,
	}
	return
}
