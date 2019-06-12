package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type DeviceAction struct {
	table *db.DeviceActions
	db    *gorm.DB
}

func GetDeviceActionAdaptor(d *gorm.DB) *DeviceAction {
	return &DeviceAction{
		table: &db.DeviceActions{Db: d},
		db:    d,
	}
}

func (n *DeviceAction) Add(device *m.DeviceAction) (id int64, err error) {

	dbDeviceAction := n.toDb(device)
	if id, err = n.table.Add(dbDeviceAction); err != nil {
		return
	}

	return
}

func (n *DeviceAction) GetById(actionId int64) (device *m.DeviceAction, err error) {

	var dbDeviceAction *db.DeviceAction
	if dbDeviceAction, err = n.table.GetById(actionId); err != nil {
		return
	}

	device = n.fromDb(dbDeviceAction)

	return
}

func (n *DeviceAction) GetByDeviceId(deviceId int64) (actions []*m.DeviceAction, err error) {

	var dbDeviceActions []*db.DeviceAction
	if dbDeviceActions, err = n.table.GetByDeviceId(deviceId); err != nil {
		return
	}

	actions = make([]*m.DeviceAction, 0)
	for _, dbActino := range dbDeviceActions {
		action := n.fromDb(dbActino)
		actions = append(actions, action)
	}

	return
}

func (n *DeviceAction) Update(device *m.DeviceAction) (err error) {
	dbDeviceAction := n.toDb(device)
	err = n.table.Update(dbDeviceAction)
	return
}

func (n *DeviceAction) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

func (n *DeviceAction) List(limit, offset int64, orderBy, sort string) (list []*m.DeviceAction, total int64, err error) {
	var dbList []*db.DeviceAction
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.DeviceAction, 0)
	for _, dbDeviceAction := range dbList {
		device := n.fromDb(dbDeviceAction)
		list = append(list, device)
	}

	return
}

func (n *DeviceAction) Search(query string, limit, offset int) (list []*m.DeviceAction, total int64, err error) {
	var dbList []*db.DeviceAction
	if dbList, total, err = n.table.Search(query, limit, offset); err != nil {
		return
	}

	list = make([]*m.DeviceAction, 0)
	for _, dbDeviceAction := range dbList {
		ver := n.fromDb(dbDeviceAction)
		list = append(list, ver)
	}

	return
}

func (n *DeviceAction) fromDb(dbDeviceAction *db.DeviceAction) (device *m.DeviceAction) {

	device = &m.DeviceAction{
		Id:          dbDeviceAction.Id,
		Name:        dbDeviceAction.Name,
		Description: dbDeviceAction.Description,
		DeviceId:    dbDeviceAction.DeviceId,
		ScriptId:    dbDeviceAction.ScriptId,
		CreatedAt:   dbDeviceAction.CreatedAt,
		UpdatedAt:   dbDeviceAction.UpdatedAt,
	}

	// device
	if dbDeviceAction.Device != nil {
		deviceAdaptor := GetDeviceAdaptor(n.db)
		device.Device = deviceAdaptor.fromDb(dbDeviceAction.Device)
	}

	// script
	if dbDeviceAction.Script != nil {
		scriptADaptor := GetScriptAdaptor(n.db)
		device.Script, _ = scriptADaptor.fromDb(dbDeviceAction.Script)
	}

	return
}

func (n *DeviceAction) toDb(device *m.DeviceAction) (dbDeviceAction *db.DeviceAction) {
	dbDeviceAction = &db.DeviceAction{
		Id:          device.Id,
		Name:        device.Name,
		Description: device.Description,
		DeviceId:    device.DeviceId,
		ScriptId:    device.ScriptId,
	}
	return
}
