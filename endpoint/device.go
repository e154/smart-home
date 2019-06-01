package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

type DeviceEndpoint struct {
	*CommonEndpoint
}

func NewDeviceEndpoint(common *CommonEndpoint) *DeviceEndpoint {
	return &DeviceEndpoint{
		CommonEndpoint: common,
	}
}

func (d *DeviceEndpoint) Add(params *m.Device) (device *m.Device, errs []*validation.Error, err error) {

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

func (d *DeviceEndpoint) GetById(deviceId int64) (device *m.Device, err error) {
	device, err = d.adaptors.Device.GetById(deviceId)
	return
}

func (d *DeviceEndpoint) Update(device *m.Device) (result *m.Device, errs []*validation.Error, err error) {

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

func (d *DeviceEndpoint) GetList(limit, offset int64, order, sortBy string) (devices []*m.Device, total int64, err error) {

	devices, total, err = d.adaptors.Device.List(limit, offset, order, sortBy)

	return
}

func (d *DeviceEndpoint) Delete(deviceId int64) (err error) {

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

func (d *DeviceEndpoint) Search(query string, limit, offset int) (devices []*m.Device, total int64, err error) {

	devices, total, err = d.adaptors.Device.Search(query, limit, offset)

	return
}
