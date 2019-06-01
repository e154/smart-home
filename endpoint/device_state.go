package endpoint

import (
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
	"errors"
)

type DeviceStateCommand struct {
	*CommonCommand
}

func NewDeviceStateCommand(common *CommonCommand) *DeviceStateCommand {
	return &DeviceStateCommand{
		CommonCommand: common,
	}
}

func (d *DeviceStateCommand) Add(params *m.DeviceState) (state *m.DeviceState, id int64, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = d.adaptors.DeviceState.Add(params); err != nil {
		return
	}

	state, err = d.GetById(id)

	return
}

func (d *DeviceStateCommand) GetById(id int64) (device *m.DeviceState, err error) {

	device, err = d.adaptors.DeviceState.GetById(id)

	return
}

func (d *DeviceStateCommand) Update(params *m.DeviceState) (state *m.DeviceState, errs []*validation.Error, err error) {

	if state, err = d.adaptors.DeviceState.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&state, &params)

	// validation
	_, errs = state.Valid()
	if len(errs) > 0 {
		return
	}

	if err = d.adaptors.DeviceState.Update(state); err != nil {
		return
	}

	state, err = d.adaptors.DeviceState.GetById(state.Id)

	return
}

func (d *DeviceStateCommand) Delete(id int64) (err error) {

	if id == 0 {
		err = errors.New("action id is null")
		return
	}

	var device *m.DeviceState
	if device, err = d.adaptors.DeviceState.GetById(id); err != nil {
		return
	}

	err = d.adaptors.DeviceState.Delete(device.Id)

	return
}

func (d *DeviceStateCommand) GetList(deviceId int64) (actions []*m.DeviceState, err error) {

	actions, err = d.adaptors.DeviceState.GetByDeviceId(deviceId)

	return
}
