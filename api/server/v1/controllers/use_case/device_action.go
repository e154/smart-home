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

func AddDeviceAction(params models.NewDeviceAction, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	action := &m.DeviceAction{}
	copier.Copy(&action, &params)

	if params.Device != nil && params.Device.Id != 0 {
		action.DeviceId = params.Device.Id
	}
	if params.Script != nil && params.Script.Id != 0 {
		action.ScriptId = params.Script.Id
	}

	// validation
	ok, errs = action.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.DeviceAction.Add(action); err != nil {
		return
	}

	action.Id = id

	return
}

func UpdateDeviceAction(params models.UpdateDeviceAction, id int64, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	action := &m.DeviceAction{}
	copier.Copy(&action, &params)

	action.Id = id

	if params.Device != nil && params.Device.Id != 0 {
		action.DeviceId = params.Device.Id
	}
	if params.Script != nil && params.Script.Id != 0 {
		action.ScriptId = params.Script.Id
	}


	// validation
	ok, errs = action.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.DeviceAction.Update(action); err != nil {
		return
	}

	return
}

func GetDeviceActionById(id int64, adaptors *adaptors.Adaptors) (device *m.DeviceAction, err error) {

	device, err = adaptors.DeviceAction.GetById(id)

	return
}

func DeleteDeviceActionById(id int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if id == 0 {
		err = errors.New("action id is null")
		return
	}

	var device *m.DeviceAction
	if device, err = adaptors.DeviceAction.GetById(id); err != nil {
		return
	}

	err = adaptors.DeviceAction.Delete(device.Id)

	return
}

func GetDeviceActionList(deviceId int64, adaptors *adaptors.Adaptors) (actions []*m.DeviceAction, err error) {

	actions, err = adaptors.DeviceAction.GetByDeviceId(deviceId)

	return
}

func SearchDeviceAction(query string, limit, offset int, adaptors *adaptors.Adaptors) (actions []*m.DeviceAction, total int64, err error) {

	actions, total, err = adaptors.DeviceAction.Search(query, limit, offset)

	return
}
