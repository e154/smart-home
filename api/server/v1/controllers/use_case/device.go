package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

func AddDevice(device *m.Device, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation
	ok, errs = device.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.Device.Add(device); err != nil {
		return
	}

	device.Id = id

	return
}

func GetDeviceById(deviceId int64, adaptors *adaptors.Adaptors) (device *m.Device, err error) {

	device, err = adaptors.Device.GetById(deviceId)

	return
}

func UpdateDevice(device *m.Device, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = device.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.Device.Update(device); err != nil {
		return
	}

	return
}

func GetDeviceList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Device, total int64, err error) {

	items, total, err = adaptors.Device.List(limit, offset, order, sortBy)

	return
}

func DeleteDeviceById(deviceId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if deviceId == 0 {
		err = errors.New("device id is null")
		return
	}

	var device *m.Device
	if device, err = adaptors.Device.GetById(deviceId); err != nil {
		return
	}

	err = adaptors.Device.Delete(device.Id)

	return
}

