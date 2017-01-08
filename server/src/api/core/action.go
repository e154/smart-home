package core

import (
	"../models"
	"../scripts"
	"../stream"
	"encoding/json"
	"reflect"
	"../log"
	"github.com/astaxie/beego/orm"
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

// ------------------------------------------------
// stream
// ------------------------------------------------
func streamDoAction(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	var device_action_id, device_id float64
	var err error

	if device_action_id, ok = v["action_id"].(float64); !ok {
		log.Warn("bad device_action_id param")
		return
	}

	if device_id, ok = v["device_id"].(float64); !ok {
		log.Warn("bad device_id param")
		return
	}

	var device_action *models.DeviceAction
	if device_action, err = models.GetDeviceActionById(int64(device_action_id)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	var device *models.Device
	if device, err = models.GetDeviceById(int64(device_id)); err != nil {
		client.Notify("error", err.Error())
		return
	}

	// get device
	// ------------------------------------------------
	var devices []*models.Device
	if device.Address != nil {
		devices = append(devices, device)
	} else {
		// значит тут группа устройств
		var childs []*models.Device
		if childs, _, err = device_action.Device.GetChilds(); err != nil {
			return
		}

		for _, child := range childs {
			if child.Address == nil || child.Status != "enabled" {
				continue
			}

			device := &models.Device{}
			*device = *device_action.Device
			device.Id = child.Id
			device.Name = child.Name
			device.Address = new(int)
			*device.Address = *child.Address
			device.Device = &models.Device{Id:int64(device_id)}
			device.Tty = child.Tty
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	if device_action.Device.Node == nil {
		client.Notify("error", "device node is nil")
		return
	}

	nodes := corePtr.GetNodes()
	var node *models.Node
	if _, ok := nodes[device_action.Device.Node.Id]; ok {
		node = nodes[device_action.Device.Node.Id]
	} else {
		// autoload nodes
		if node, err = models.GetNodeById(device_action.Device.Node.Id); err != nil {
			client.Notify("error", err.Error())
			return
		}

		if err = CorePtr().AddNode(node); err != nil {
			client.Notify("error", err.Error())
			return
		}
	}

	// get script
	// ------------------------------------------------
	o := orm.NewOrm()
	if _, err = o.LoadRelated(device_action, "Script"); err != nil {
		client.Notify("error", err.Error())
		return
	}

	for _, device := range devices {
		var action *Action
		if action, err = NewAction(device, device_action.Script, node); err != nil {
			log.Error(err.Error())
			client.Notify("error", err.Error())
			continue
		}

		action.Do()
		//body, _ := action.Do()
		//client.Notify("success", body)
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send(string(msg))
}
