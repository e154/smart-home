package use_case

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
	"errors"
	"github.com/e154/smart-home/common"
)

func AddDeviceState(params models.NewDeviceState, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	state := &m.DeviceState{}
	copier.Copy(&state, &params)

	if params.Device != nil && params.Device.Id != 0 {
		state.DeviceId = params.Device.Id
	}

	// validation
	ok, errs = state.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.DeviceState.Add(state); err != nil {
		return
	}

	state.Id = id

	return
}

func UpdateDeviceState(params models.UpdateDeviceState, id int64, adaptors *adaptors.Adaptors, core *core.Core) (result *models.DeviceState, errs []*validation.Error, err error) {

	var state *m.DeviceState
	if state, err = adaptors.DeviceState.GetById(id); err != nil {
		return
	}

	copier.Copy(&state, &params)

	if params.Device != nil && params.Device.Id != 0 {
		state.DeviceId = params.Device.Id
	}

	// validation
	_, errs = state.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.DeviceState.Update(state); err != nil {
		return
	}


	if state, err = adaptors.DeviceState.GetById(id); err != nil {
		return
	}

	result = &models.DeviceState{}
	err = common.Copy(&result, &state)

	return
}

func GetDeviceStateById(id int64, adaptors *adaptors.Adaptors) (device *m.DeviceState, err error) {

	device, err = adaptors.DeviceState.GetById(id)

	return
}

func DeleteDeviceStateById(id int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if id == 0 {
		err = errors.New("action id is null")
		return
	}

	var device *m.DeviceState
	if device, err = adaptors.DeviceState.GetById(id); err != nil {
		return
	}

	err = adaptors.DeviceState.Delete(device.Id)

	return
}

func GetDeviceStateList(deviceId int64, adaptors *adaptors.Adaptors) (actions []*m.DeviceState, err error) {

	actions, err = adaptors.DeviceState.GetByDeviceId(deviceId)

	return
}
