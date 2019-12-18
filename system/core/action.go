package core

import (
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
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

	a.ScriptEngine.PushGlobalProxy("device", &DeviceBind{
		model: a.Device,
		node:  a.Node,
	})

	if a.flow != nil {
		a.ScriptEngine.PushGlobalProxy("flow", &FlowBind{flow: a.flow})
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
	a.ScriptEngine.PushGlobalProxy("message", a.Message)

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
