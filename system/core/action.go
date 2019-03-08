package core

import (
	"github.com/e154/smart-home/system/scripts"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/common"
)

type Action struct {
	Device        *m.Device
	Node          *Node
	ScriptEngine  *scripts.Engine
	Message       *Message
	flow          *Flow
	scriptService *scripts.ScriptService
	script        *m.Script
}

func NewAction(device *m.Device,
	script *m.Script,
	node *Node,
	flow *Flow,
	scriptService *scripts.ScriptService) (action *Action, err error) {

	action = &Action{
		Device:        device,
		Node:          node,
		flow:          flow,
		scriptService: scriptService,
		script:        script,
	}

	err = action.NewScript()

	return
}

func (a *Action) Do() (res string, err error) {

	a.ScriptEngine.PushStruct("device", &DeviceBind{
		model: a.Device,
		node:  a.Node,
	})

	a.ScriptEngine.PushStruct("flow", &FlowBind{flow: a.flow})

	a.Message.Clear()
	/*res,*/ err = a.ScriptEngine.EvalString(a.script.Compiled)
	//a.Message.SetVar("result", res)
	return
}

func (a *Action) NewScript() (err error) {

	if a.flow != nil {
		if a.ScriptEngine, err = a.flow.NewScript(); err != nil {
			return
		}

	} else {
		model := &m.Script{
			Lang: ScriptLangJavascript,
		}
		if a.ScriptEngine, err = a.scriptService.NewEngine(model); err != nil {
			return
		}
	}

	a.Message = NewMessage()
	a.ScriptEngine.PushStruct("message", a.Message)

	// bind
	javascript := a.ScriptEngine.Get().(*scripts.Javascript)
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

func (a *Action) GetDevice() *m.Device {
	return a.Device
}

func (a *Action) GetNode() *Node {
	return a.Node
}

// ------------------------------------------------
// stream
// ------------------------------------------------
//func streamDoAction(client *stream.Client, value interface{}) {
//v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
//if !ok {
//	return
//}
//
//var device_action_id, device_id float64
//var err error
//
//if device_action_id, ok = v["action_id"].(float64); !ok {
//	log.Warning("bad device_action_id param")
//	return
//}
//
//if device_id, ok = v["device_id"].(float64); !ok {
//	log.Warning("bad device_id param")
//	return
//}

//	var device_action *m.DeviceAction
//	if device_action, err = m.GetDeviceActionById(int64(device_action_id)); err != nil {
//		client.Notify("error", err.Error())
//		return
//	}
//
//	var device *m.Device
//	if device, err = m.GetDeviceById(int64(device_id)); err != nil {
//		client.Notify("error", err.Error())
//		return
//	}
//
//	// get device
//	// ------------------------------------------------
//	var devices []*m.Device
//	if device.Address != nil {
//		devices = append(devices, device)
//	} else {
//		// значит тут группа устройств
//		var childs []*m.Device
//		if childs, _, err = device_action.Device.GetChilds(); err != nil {
//			return
//		}
//
//		for _, child := range childs {
//			if child.Address == nil || child.Status != "enabled" {
//				continue
//			}
//
//			device := &m.Device{}
//			*device = *device_action.Device
//			device.Id = child.Id
//			device.Name = child.Name
//			device.Address = new(int)
//			*device.Address = *child.Address
//			device.Device = &m.Device{Id:int64(device_id)}
//			device.Tty = child.Tty
//			device.Sleep = device_action.Device.Sleep
//			devices = append(devices, device)
//		}
//	}
//
//	// get node
//	// ------------------------------------------------
//	if device_action.Device.NodeModel == nil {
//		client.Notify("error", "device node is nil")
//		return
//	}
//
//	nodes := corePtr.GetNodes()
//	var node *NodeModel
//	if _, ok := nodes[device_action.Device.NodeModel.Id]; ok {
//		node = nodes[device_action.Device.NodeModel.Id]
//	} else {
//		// autoload nodes
//		var model_node *m.NodeModel
//		model_node, err = models.GetNodeById(device_action.Device.NodeModel.Id)
//		if err != nil {
//			client.Notify("error", err.Error())
//			return
//		}
//
//		if err = CorePtr().AddNode(NewNode(model_node)); err != nil {
//			client.Notify("error", err.Error())
//			return
//		}
//	}
//
//	// get script
//	// ------------------------------------------------
//	o := orm.NewOrm()
//	if _, err = o.LoadRelated(device_action, "Script"); err != nil {
//		client.Notify("error", err.Error())
//		return
//	}
//
//	for _, device := range devices {
//		var action *Action
//		if action, err = NewAction(device, device_action.Script, node, nil); err != nil {
//			log.Error(err.Error())
//			client.Notify("error", err.Error())
//			continue
//		}
//
//		action.Do()
//		//body, _ := action.Do()
//		//client.Notify("success", body)
//	}
//
//	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
//	client.Send <- msg
//}
