package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/t-tiger/gorm-bulk-insert"
)

type MapDeviceAction struct {
	table *db.MapDeviceActions
	db    *gorm.DB
}

func GetMapDeviceActionAdaptor(d *gorm.DB) *MapDeviceAction {
	return &MapDeviceAction{
		table: &db.MapDeviceActions{Db: d},
		db:    d,
	}
}

func (n *MapDeviceAction) Add(ver *m.MapDeviceAction) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapDeviceAction) AddMultiple(items []*m.MapDeviceAction) (err error) {

	insertRecords := make([]interface{}, 0)
	for _, ver := range items {
		dbVer := n.toDb(ver)
		insertRecords = append(insertRecords, dbVer)
	}

	err = gormbulk.BulkInsert(n.db, insertRecords, 3000)

	return
}

func (n *MapDeviceAction) fromDb(dbVer *db.MapDeviceAction) (ver *m.MapDeviceAction) {
	ver = &m.MapDeviceAction{
		Id:             dbVer.Id,
		MapDeviceId:    dbVer.MapDeviceId,
		ImageId:        dbVer.ImageId,
		Type:           dbVer.Type,
		DeviceActionId: dbVer.DeviceActionId,
		CreatedAt:      dbVer.CreatedAt,
		UpdatedAt:      dbVer.UpdatedAt,
	}

	return
}

func (n *MapDeviceAction) toDb(ver *m.MapDeviceAction) (dbVer *db.MapDeviceAction) {
	dbVer = &db.MapDeviceAction{
		Id:             ver.Id,
		MapDeviceId:    ver.MapDeviceId,
		ImageId:        ver.ImageId,
		Type:           ver.Type,
		DeviceActionId: ver.DeviceActionId,
	}
	return
}
