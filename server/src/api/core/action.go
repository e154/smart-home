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
	//TODO refactor message system
	action.Message = NewMessage()
	action.Message.Device = device
	action.Message.Node = action.Node
	action.Message.Device_state = func(state string) {
		action.SetDeviceState(state)
	}

	action.Script.PushStruct("message", action.Message)

	return
}

type Action struct {
	Device  *models.Device
	Node    *models.Node
	Script  *scripts.Engine
	Message *Message
}

func (a *Action) Do() (res string, err error) {
	//TODO refactor message system
	a.Message.Error = ""
	a.Message.Data = make(map[string]interface{})
	res, err = a.Script.Do()
	a.Message.Data["result"] = res
	return
}

func (a *Action) SetDeviceState(_state string) {
	for _, state := range a.Device.States {
		if state.SystemName == _state {
			CorePtr().SetDeviceState(a.Device.Id, state)
			break
		}
	}
}