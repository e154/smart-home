package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type MapDevice struct {
	table *db.MapDevices
	db    *gorm.DB
}

func GetMapDeviceAdaptor(d *gorm.DB) *MapDevice {
	return &MapDevice{
		table: &db.MapDevices{Db: d},
		db:    d,
	}
}

func (n *MapDevice) Add(ver *m.MapDevice) (id int64, err error) {

	dbVer := n.toDb(ver)
	if id, err = n.table.Add(dbVer); err != nil {
		return
	}

	return
}

func (n *MapDevice) GetById(mapId int64) (ver *m.MapDevice, err error) {

	var dbVer *db.MapDevice
	if dbVer, err = n.table.GetById(mapId); err != nil {
		return
	}

	ver = n.fromDb(dbVer)

	return
}

func (n *MapDevice) Delete(mapId int64) (err error) {
	err = n.table.Delete(mapId)
	return
}

func (n *MapDevice) fromDb(dbVer *db.MapDevice) (ver *m.MapDevice) {
	ver = &m.MapDevice{
		Id:         dbVer.Id,
		SystemName: dbVer.SystemName,
		DeviceId:   dbVer.DeviceId,
		ImageId:    dbVer.ImageId,
		Actions:    make([]*m.MapDeviceAction, 0),
		States:     make([]*m.MapDeviceState, 0),
		CreatedAt:  dbVer.CreatedAt,
		UpdatedAt:  dbVer.UpdatedAt,
	}

	// actions
	mapDeviceActionAdaptor := GetMapDeviceActionAdaptor(n.db)
	for _, dbAction := range dbVer.Actions {
		action := mapDeviceActionAdaptor.fromDb(dbAction)
		ver.Actions = append(ver.Actions, action)
	}

	// states
	mapDeviceStateAdaptor := GetMapDeviceStateAdaptor(n.db)
	for _, dbState := range dbVer.States {
		state := mapDeviceStateAdaptor.fromDb(dbState)
		ver.States = append(ver.States, state)
	}

	// device
	if dbVer.Device != nil {
		deviceAdaptor := GetDeviceAdaptor(n.db)
		ver.Device = deviceAdaptor.fromDb(dbVer.Device)
	}

	// image
	if dbVer.Image != nil {
		imageAdaptor := GetImageAdaptor(n.db)
		ver.Image = imageAdaptor.fromDb(dbVer.Image)
	}

	return
}

func (n *MapDevice) toDb(ver *m.MapDevice) (dbVer *db.MapDevice) {
	dbVer = &db.MapDevice{
		Id:         ver.Id,
		SystemName: ver.SystemName,
		DeviceId:   ver.DeviceId,
		ImageId:    ver.ImageId,
	}
	return
}
