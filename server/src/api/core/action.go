package core

import (
	"../models"
	"../scripts"
	r "../../lib/rpc"
)

func NewAction(device *models.Device, device_action *models.DeviceAction, node *models.Node) (action *Action, err error) {

	action = &Action{
		Device: 	device,
		DeviceAction:	device_action,
		Node:		node,
	}

	// add script
	// ------------------------------------------------
	if device_action.Script == nil {
		return
	}

	if action.Script, err = scripts.New(device_action.Script); err != nil {
		return
	}

	action.Script.PushStruct("device", device)
	action.Script.PushStruct("node", action.Node)
	action.Script.PushStruct("request", &r.Request{})

	action.Script.PushFunction("modbus_send", func(args *r.Request) (result r.Result) {
		if err := action.Node.ModbusSend(args, &result); err != nil {
			result.Error = err.Error()
		}

		return
	})

	return
}

type Action struct {
	Device		*models.Device
	DeviceAction	*models.DeviceAction
	Node		*models.Node
	Script		*scripts.Engine
}

func (a *Action) Do() (string, error) {
	return a.Script.Do()
}