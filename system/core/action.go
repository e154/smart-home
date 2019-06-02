package core

import (
	"github.com/e154/smart-home/system/scripts"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/common"
	"reflect"
	"github.com/e154/smart-home/system/stream"
	"encoding/json"
	"fmt"
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

	if a.flow != nil {
		a.ScriptEngine.PushStruct("flow", &FlowBind{flow: a.flow})
	}

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
func streamDoAction(core *Core) func(client *stream.Client, value interface{}) {

	return func(client *stream.Client, value interface{}) {

		v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
		if !ok {
			return
		}

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
		if device, err = core.adaptors.Device.GetById(int64(deviceId)); err != nil {
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
		var node *Node
		if device.Node != nil {
			node = core.GetNodeById(device.Node.Id)
		} else {
			client.Notify("error", "node in device is nil")
			return
		}

		if node == nil {
			client.Notify("error", fmt.Sprintf("node id(%v) not found", node.Id))
			return
		}

		// action
		var action *Action
		if action, err = NewAction(device, deviceAction.Script, node, nil, core.scripts); err != nil {
			client.Notify("error", err.Error())
			return
		}

		// do action
		if _, err = action.Do(); err != nil {
			log.Error(err.Error())
		}

		msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "status": "ok"})
		client.Send <- msg
	}
}
