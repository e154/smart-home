package controllers

import (
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/stream"
)

type ControllerAction struct {
	*ControllerCommon
}

func NewControllerAction(common *ControllerCommon) *ControllerAction {
	return &ControllerAction{
		ControllerCommon: common,
	}
}

func (c *ControllerAction) Start() {
	c.stream.Subscribe("do.action", c.DoAction)
}


func (c *ControllerAction) Stop() {
	c.stream.UnSubscribe("do.action")
}

// Stream
func (c *ControllerAction) DoAction(client *stream.Client, message stream.Message) {

	v := message.Payload
	var ok bool

	var deviceActionId, deviceId float64
	var err error

	if deviceActionId, ok = v["action_id"].(float64); !ok {
		log.Warning("bad device_action_id param")
		return
	}

	if deviceId, ok = v["device_id"].(float64); !ok {
		log.Warning("bad device_id param")
		return
	}

	// device
	var device *m.Device
	if device, err = c.adaptors.Device.GetById(int64(deviceId)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	// device action
	var deviceAction *m.DeviceAction
	for _, action := range device.Actions {
		if action.Id == int64(deviceActionId) {
			deviceAction = action
		}
	}

	if deviceAction == nil {
		client.Notify("error", fmt.Sprintf("device action id(%v) not found", deviceActionId))
		return
	}

	// node
	var node *core.Node
	if device.Node != nil {
		node = c.core.GetNodeById(device.Node.Id)
	} else {
		client.Notify("error", "node in device is nil")
		return
	}

	if node == nil {
		client.Notify("error", fmt.Sprintf("node id(%v) not found", node.Id))
		return
	}

	// action
	var action *core.Action
	if action, err = core.NewAction(device, deviceAction.Script, node, nil, c.scripts); err != nil {
		client.Notify("error", err.Error())
		return
	}

	// do action
	if _, err = action.Do(); err != nil {
		c.Err(client, message, err)
	}

	client.Send <- message.Success().Pack()
}
