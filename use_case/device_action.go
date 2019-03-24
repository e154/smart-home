package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

type DeviceActionCommand struct {
	*CommonCommand
}

func NewDeviceActionCommand(common *CommonCommand) *DeviceActionCommand {
	return &DeviceActionCommand{
		CommonCommand: common,
	}
}

func (d *DeviceActionCommand) Add(params *m.DeviceAction) (action *m.DeviceAction, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = d.adaptors.DeviceAction.Add(params); err != nil {
		return
	}

	action, err = d.adaptors.DeviceAction.GetById(id)

	return
}

func (d *DeviceActionCommand) Update(params *m.DeviceAction) (result *m.DeviceAction, errs []*validation.Error, err error) {

	var action *m.DeviceAction
	if action, err = d.adaptors.DeviceAction.GetById(params.Id); err != nil {
		return
	}

	// validation
	_, errs = action.Valid()
	if len(errs) > 0 {
		return
	}

	if err = d.adaptors.DeviceAction.Update(action); err != nil {
		return
	}

	action, err = d.adaptors.DeviceAction.GetById(params.Id)

	return
}

func (d *DeviceActionCommand) GetById(id int64) (device *m.DeviceAction, err error) {

	device, err = d.adaptors.DeviceAction.GetById(id)

	return
}

func (d *DeviceActionCommand) Delete(id int64) (err error) {

	if id == 0 {
		err = errors.New("action id is null")
		return
	}

	var device *m.DeviceAction
	if device, err = d.adaptors.DeviceAction.GetById(id); err != nil {
		return
	}

	err = d.adaptors.DeviceAction.Delete(device.Id)

	return
}

func (d *DeviceActionCommand) GetList(deviceId int64) (actions []*m.DeviceAction, err error) {

	actions, err = d.adaptors.DeviceAction.GetByDeviceId(deviceId)

	return
}

func (d *DeviceActionCommand) Search(query string, limit, offset int) (actions []*m.DeviceAction, total int64, err error) {

	actions, total, err = d.adaptors.DeviceAction.Search(query, limit, offset)

	return
}
