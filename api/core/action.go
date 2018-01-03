package core

import (
	"reflect"
	"encoding/json"
	"github.com/e154/smart-home/api/models"
	"github.com/e154/smart-home/api/scripts"
	"github.com/e154/smart-home/api/stream"
	"github.com/e154/smart-home/api/log"
	"github.com/astaxie/beego/orm"
)

func NewAction(device *models.Device, script *models.Script, node *Node, flow *Flow) (action *Action, err error) {

	action = &Action{
		Device: 	device,
		Node:		node,
	}

	if err = action.Device.GetInheritedData(); err != nil {
		return
	}

	err = action.NewScript(flow, script)

	return
}

type Action struct {
	Device  *models.Device
	Node    *Node
	Script  *scripts.Engine
	Message *Message
}

func (a *Action) Do() (res string, err error) {
	//TODO refactor message system
	a.Message.Clear()
	res, err = a.Script.Do()
	a.Message.SetVar("result", res)
	return
}

func (a *Action) NewScript(flow *Flow, script *models.Script) (err error) {

	if flow != nil {
		if a.Script, err = flow.NewScript(script); err != nil {
			return
		}
	} else {
		if a.Script, err = scripts.New(script); err != nil {
			return
		}
	}

	a.Message = NewMessage()
	a.Script.PushStruct("message", a.Message)

	// bind
	javascript := a.Script.Get().(*scripts.Javascript)
	ctx := javascript.Ctx()
	if b := ctx.GetGlobalString("IC"); !b {
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() *ActionBind {
		return &ActionBind{action: a}
	})
	ctx.PutPropString(-3, "Action")
	ctx.Pop()

	return nil
}

func (a *Action) GetDevice() *models.Device {
	return a.Device
}

func (a *Action) GetNode() *Node {
	return a.Node
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
			device.Sleep = device_action.Device.Sleep
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
	var node *Node
	if _, ok := nodes[device_action.Device.Node.Id]; ok {
		node = nodes[device_action.Device.Node.Id]
	} else {
		// autoload nodes
		var model_node *models.Node
		model_node, err = models.GetNodeById(device_action.Device.Node.Id)
		if err != nil {
			client.Notify("error", err.Error())
			return
		}

		if err = CorePtr().AddNode(NewNode(model_node)); err != nil {
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
		if action, err = NewAction(device, device_action.Script, node, nil); err != nil {
			log.Error(err.Error())
			client.Notify("error", err.Error())
			continue
		}

		action.Do()
		//body, _ := action.Do()
		//client.Notify("success", body)
	}

	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
	client.Send <- msg
}
