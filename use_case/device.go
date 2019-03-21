package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

type DeviceCommand struct {
	*CommonCommand
}

func NewDeviceCommand(common *CommonCommand) *DeviceCommand {
	return &DeviceCommand{
		CommonCommand: common,
	}
}

func (d *DeviceCommand) Add(params *m.Device) (device *m.Device, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = d.adaptors.Device.Add(params); err != nil {
		return
	}

	device, err = d.adaptors.Device.GetById(id)

	return
}

func (d *DeviceCommand) GetById(deviceId int64) (device *m.Device, err error) {
	device, err = d.adaptors.Device.GetById(deviceId)
	return
}

func (d *DeviceCommand) Update(device *m.Device) (result *m.Device, errs []*validation.Error, err error) {

	_, errs = device.Valid()
	if len(errs) > 0 {
		return
	}

	if err = d.adaptors.Device.Update(device); err != nil {
		return
	}

	result, err = d.adaptors.Device.GetById(device.Id)

	return
}

func (d *DeviceCommand) GetList(limit, offset int64, order, sortBy string) (devices []*m.Device, total int64, err error) {

	devices, total, err = d.adaptors.Device.List(limit, offset, order, sortBy)

	return
}

func (d *DeviceCommand) Delete(deviceId int64) (err error) {

	if deviceId == 0 {
		err = errors.New("device id is null")
		return
	}

	var device *m.Device
	if device, err = d.adaptors.Device.GetById(deviceId); err != nil {
		return
	}

	err = d.adaptors.Device.Delete(device.Id)

	return
}

func (d *DeviceCommand) Search(query string, limit, offset int) (devices []*m.Device, total int64, err error) {

	devices, total, err = d.adaptors.Device.Search(query, limit, offset)

	return
}
