package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	"errors"
	"encoding/json"
	"github.com/e154/smart-home/common"
)

func AddDevice(params models.NewDevice, adaptors *adaptors.Adaptors, core *core.Core) (result *models.Device, errs []*validation.Error, err error) {

	var properties []byte
	if properties, err = json.Marshal(params.Properties); err != nil {
		return
	}

	device := &m.Device{
		Name:        params.Name,
		Description: params.Description,
		Properties:  properties,
		Status:      params.Status,
		Type:        common.DeviceType(params.Type),
	}

	if params.Device != nil && params.Device.Id != 0 {
		device.Device = &m.Device{Id: params.Device.Id}
	}

	if params.Node != nil && params.Node.Id != 0 {
		device.Node = &m.Node{Id: params.Node.Id}
	}

	//device.SetPropertiesFromMap(params.Properties)

	// validation
	_, errs = device.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = adaptors.Device.Add(device); err != nil {
		return
	}

	if device, err = adaptors.Device.GetById(id); err != nil {
		return
	}

	result = &models.Device{}
	err = common.Copy(&result, &device)


	return
}

func GetDeviceById(deviceId int64, adaptors *adaptors.Adaptors) (result *models.Device, err error) {

	var device *m.Device
	if device, err = adaptors.Device.GetById(deviceId); err != nil {
		return
	}

	result = &models.Device{}
	err = common.Copy(&result, &device, common.JsonEngine)

	return
}

func UpdateDevice(params models.UpdateDevice, id int64, adaptors *adaptors.Adaptors, core *core.Core) (result *models.Device, errs []*validation.Error, err error) {

	var properties []byte
	if properties, err = json.Marshal(params.Properties); err != nil {
		return
	}

	device := &m.Device{
		Id:          id,
		Name:        params.Name,
		Description: params.Description,
		Properties:  properties,
		Status:      params.Status,
		Type:        common.DeviceType(params.Type),
	}

	//device.SetPropertiesFromMap(params.Properties)

	// validation
	_, errs = device.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Device.Update(device); err != nil {
		return
	}

	if device, err = adaptors.Device.GetById(id); err != nil {
		return
	}

	result = &models.Device{}
	err = common.Copy(&result, &device)

	return
}

func GetDeviceList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.Device, total int64, err error) {

	var devices []*m.Device
	if devices, total, err = adaptors.Device.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.Device, 0)
	err = common.Copy(&result, &devices, common.JsonEngine)

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

func SearchDevice(query string, limit, offset int, adaptors *adaptors.Adaptors) (result []*models.DeviceShort, total int64, err error) {

	var devices []*m.Device
	if devices, total, err = adaptors.Device.Search(query, limit, offset); err != nil {
		return
	}

	result = make([]*models.DeviceShort, 0)
	err = common.Copy(&result, &devices, common.JsonEngine)

	return
}
