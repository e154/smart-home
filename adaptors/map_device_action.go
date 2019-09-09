package adaptors

import (
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/gorm"
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

	for _, ver := range items {
		if ver.Image == nil {
			continue
		}
		dbVer := n.toDb(ver)
		if _, err = n.table.Add(dbVer); err != nil {
			return
		}
	}

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

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	// actions
	if dbVer.DeviceAction != nil {
		deviceActionAdaptor := GetDeviceActionAdaptor(n.db)
		ver.DeviceAction = deviceActionAdaptor.fromDb(dbVer.DeviceAction)
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
	if ver.DeviceAction != nil && ver.DeviceAction.Id != 0 {
		dbVer.DeviceActionId = ver.DeviceAction.Id
	}
	if ver.Image != nil && ver.Image.Id != 0 {
		dbVer.ImageId = ver.Image.Id
	}
	return
}
