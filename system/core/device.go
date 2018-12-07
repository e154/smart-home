package core

import (
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/models/devices"
	"encoding/json"
)

type Device struct {
	dev  *m.Device
	node *Node
}

// run command
func (d *Device) RunCommand(name string, args []string) (result *DevCommandResponse) {

	request := &DevCommandRequest{
		Name: name,
		Args: args,
	}

	result = &DevCommandResponse{}
	data, err := json.Marshal(request)
	if err != nil {
		result.Error = err.Error()
		return
	}

	nodeResult, err := d.node.Send(d.dev, data)

	if err = json.Unmarshal(nodeResult.Response, result); err != nil {
		result.Error = err.Error()
		return
	}

	//debug.Println(nodeResult)

	result.Time = nodeResult.Time

	return
}

func (d *Device) SmartBus(command []byte) (result *DevSmartBusResponse) {

	request := &DevSmartBusRequest{
		Command: command,
	}

	result = &DevSmartBusResponse{}
	data, err := json.Marshal(request)
	if err != nil {
		result.Error = err.Error()
		return
	}

	nodeResult, err := d.node.Send(d.dev, data)

	if err = json.Unmarshal(nodeResult.Response, result); err != nil {
		result.Error = err.Error()
		return
	}

	//debug.Println(nodeResult)

	result.Time = nodeResult.Time

	return
}
