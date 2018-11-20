package adaptors

import (
	"github.com/jinzhu/gorm"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

type Device struct {
	table *db.Devices
	db    *gorm.DB
}

func GetDeviceAdaptor(d *gorm.DB) *Device {
	return &Device{
		table: &db.Devices{Db: d},
		db:    d,
	}
}

func (n *Device) Add(device *m.Device) (id int64, err error) {

	dbDevice := n.toDb(device)
	if id, err = n.table.Add(dbDevice); err != nil {
		return
	}

	return
}

func (n *Device) GetAllEnabled() (list []*m.Device, err error) {

	var dbList []*db.Device
	if dbList, err = n.table.GetAllEnabled(); err != nil {
		return
	}

	list = make([]*m.Device, 0)
	for _, dbDevice := range dbList {
		device := n.fromDb(dbDevice)
		list = append(list, device)
	}

	return
}

func (n *Device) GetById(deviceId int64) (device *m.Device, err error) {

	var dbDevice *db.Device
	if dbDevice, err = n.table.GetById(deviceId); err != nil {
		return
	}

	device = n.fromDb(dbDevice)

	return
}

func (n *Device) Update(device *m.Device) (err error) {
	dbDevice := n.toDb(device)
	err = n.table.Update(dbDevice)
	return
}

func (n *Device) Delete(deviceId int64) (err error) {
	err = n.table.Delete(deviceId)
	return
}

func (n *Device) List(limit, offset int64, orderBy, sort string) (list []*m.Device, total int64, err error) {
	var dbList []*db.Device
	if dbList, total, err = n.table.List(limit, offset, orderBy, sort); err != nil {
		return
	}

	list = make([]*m.Device, 0)
	for _, dbDevice := range dbList {
		device := n.fromDb(dbDevice)
		list = append(list, device)
	}

	return
}

func (n *Device) fromDb(dbDevice *db.Device) (device *m.Device) {
	device = &m.Device{
		Id:          dbDevice.Id,
		Name:        dbDevice.Name,
		Status:      dbDevice.Status,
		Description: dbDevice.Description,
		Type:        dbDevice.Type,
		Properties:  dbDevice.Properties,
		IsGroup:     dbDevice.DeviceId == nil,
		Actions:     make([]*m.DeviceAction, 0),
		States:      make([]*m.DeviceState, 0),
		Devices:     make([]*m.Device, 0),
		NodeId:      dbDevice.NodeId,
		CreatedAt:   dbDevice.CreatedAt,
		UpdatedAt:   dbDevice.UpdatedAt,
	}

	// actions
	deviceActionAdaptor := GetDeviceActionAdaptor(n.db)
	for _, dbAction := range dbDevice.Actions {
		action := deviceActionAdaptor.fromDb(dbAction)
		device.Actions = append(device.Actions, action)
	}

	// states
	deviceStatesAdaptor := GetDeviceStateAdaptor(n.db)
	for _, dbState := range dbDevice.States {
		state := deviceStatesAdaptor.fromDb(dbState)
		device.States = append(device.States, state)
	}

	// devices
	for _, dbDevice := range dbDevice.Devices {
		dev := n.fromDb(dbDevice)
		device.Devices = append(device.Devices, dev)
	}

	return
}

func (n *Device) toDb(device *m.Device) (dbDevice *db.Device) {
	dbDevice = &db.Device{
		Id:          device.Id,
		Name:        device.Name,
		Status:      device.Status,
		Description: device.Description,
		CreatedAt:   device.CreatedAt,
		UpdatedAt:   device.UpdatedAt,
		NodeId:      device.NodeId,
		Properties:  device.Properties,
		Type:        device.Type,
	}

	// device
	if device.Device != nil {
		dbDevice.DeviceId = &device.Device.Id
	}

	return
}
