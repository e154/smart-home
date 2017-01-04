package core

import (
	"../models"
	"../scripts"
)

func NewAction(device *models.Device, script *models.Script, node *models.Node) (action *Action, err error) {

	action = &Action{
		Device: 	device,
		Node:		node,
	}

	// add script
	// ------------------------------------------------
	if action.Script, err = scripts.New(script); err != nil {
		return
	}

	action.Device.GetInheritedData()

	message := &Message{
		Device: device,
		Node: action.Node,
	}

	action.Script.PushStruct("message", message)
	action.Script.PushFunction("set_device_state", func(state string) {
		action.SetDeviceState(state)
	})

	return
}

type Action struct {
	Device		*models.Device
	Node		*models.Node
	Script		*scripts.Engine
}

func (a *Action) Do() (string, error) {
	return a.Script.Do()
}

func (a *Action) SetDeviceState(_state string) {
	for _, state := range a.Device.States {
		if state.SystemName == _state {
			CorePtr().SetDeviceState(a.Device.Id, state)
			break
		}
	}
}