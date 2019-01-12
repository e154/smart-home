package use_case

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
	"errors"
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

func UpdateDeviceState(params models.UpdateDeviceState, id int64, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	action := &m.DeviceState{}
	copier.Copy(&action, &params)

	action.Id = id

	if params.Device != nil && params.Device.Id != 0 {
		action.DeviceId = params.Device.Id
	}

	// validation
	ok, errs = action.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.DeviceState.Update(action); err != nil {
		return
	}

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
