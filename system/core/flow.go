package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	cr "github.com/e154/smart-home/system/cron"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/uuid"
)

type Flow struct {
	Storage
	Model         *m.Flow
	workflow      *Workflow
	Connections   []*m.Connection
	FlowElements  []*FlowElement
	cursor        uuid.UUID
	Node          *Node
	quit          chan bool
	adaptors      *adaptors.Adaptors
	scriptService *scripts.ScriptService
	scriptEngine  *scripts.Engine
	Workers       map[int64]*Worker
	cron          *cr.Cron
	core          *Core
}

func NewFlow(model *m.Flow,
	workflow *Workflow,
	adaptors *adaptors.Adaptors,
	scripts *scripts.ScriptService,
	cron *cr.Cron,
	core *Core) (flow *Flow, err error) {

	flow = &Flow{
		Model:         model,
		workflow:      workflow,
		quit:          make(chan bool),
		adaptors:      adaptors,
		scriptService: scripts,
		Workers:       make(map[int64]*Worker),
		cron:          cron,
		core:          core,
	}

	flow.pull = make(map[string]interface{})

	if flow.scriptEngine, err = flow.NewScript(); err != nil {
		return
	}

	for _, element := range flow.Model.FlowElements {
		var flowElement *FlowElement
		if flowElement, err = NewFlowElement(element, flow, workflow, adaptors); err == nil {
			flow.FlowElements = append(flow.FlowElements, flowElement)
		} else {
			log.Warning(err.Error())
		}
	}

	for _, conn := range flow.Model.Connections {
		flow.Connections = append(flow.Connections, conn)
	}

	// add worker
	err = flow.InitWorkers()

	return
}

func (f *Flow) Remove() {
	//f.quit <- true
	for _, worker := range f.Workers {
		f.RemoveWorker(worker.Model)
	}
}

func (f *Flow) NewMessage(ctx context.Context) (err error) {

	// circular dependency search
	if ctx, err = f.defineCircularConnection(ctx); err != nil {
		return
	}

	var _element *FlowElement

	// find message handler
	// ------------------------------------------------
	for _, element := range f.FlowElements {
		if element.Prototype == nil {
			continue
		}

		if element.Model.PrototypeType != "MessageHandler" {
			continue
		}

		_element = element
		break
	}

	if _element == nil {
		err = errors.New("message handler not found")
		return
	}

	// ------------------------------------------------
	getNextElements := func(element *FlowElement, isScripted, isTrue bool) (elements []*FlowElement) {
		// each connections
		for _, conn := range f.Connections {
			if conn.ElementFrom != element.Model.Uuid || conn.ElementTo == element.Model.Uuid {
				continue
			}

			for _, element := range f.FlowElements {
				if conn.ElementTo != element.Model.Uuid {
					continue
				}

				if isScripted {
					if conn.Direction == "true" {
						if !isTrue {
							continue
						}
					} else if conn.Direction == "false" {
						if isTrue {
							continue
						}
					}
				}

				elements = append(elements, element)
			}
		}

		return
	}

	var message *Message

	var runElement func(*FlowElement)
	var returnMessage *Message
	runElement = func(element *FlowElement) {

		if err = ctx.Err(); err != nil {
			return
		}

		var ok, isScripted bool
		isScripted = element.ScriptEngine != nil

		childCtx, _ := context.WithCancel(ctx)
		if message != nil {
			childCtx = context.WithValue(childCtx, "msg", message)
		}

		if ok, returnMessage, err = element.Run(childCtx); err != nil {
			log.Error(err.Error())
			return
		}

		// send message to linked flow
		if element.Model.PrototypeType == "Flow" && element.Model.FlowLink != nil {
			if flow, ok := f.workflow.Flows[*element.Model.FlowLink]; ok {
				childCtx, _ := context.WithCancel(ctx)
				if err = flow.NewMessage(childCtx); err != nil {
					return
				}
			}
		}

		// copy message
		if returnMessage != nil {
			message = returnMessage
		}

		elements := getNextElements(element, isScripted, ok)
		for _, e := range elements {
			runElement(e)
		}
	}

	runElement(_element)

	return
}

func (f *Flow) loop() {
	//for {
	//	select {
	//	case <-f.quit:
	//		break
	//	}
	//}
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (f *Flow) InitWorkers() (err error) {

	for _, worker := range f.Model.Workers {
		if err = f.AddWorker(worker); err != nil {
			return
		}
	}

	return
}

func (f *Flow) AddWorker(model *m.Worker) (err error) {

	log.Infof("Add worker: \"%s\"", model.Name)

	if _, ok := f.Workers[model.Id]; ok {
		return
	}

	if len(f.FlowElements) == 0 {
		err = errors.New("no flow elements")
		return
	}

	// get device
	// ------------------------------------------------
	var devices []*m.Device
	if !model.DeviceAction.Device.IsGroup {
		devices = append(devices, model.DeviceAction.Device)
	} else {
		// значит тут группа устройств
		for _, child := range model.DeviceAction.Device.Devices {
			if child.Status != "enabled" {
				continue
			}

			//if child.Address == nil {
			//	continue
			//}

			device := &m.Device{
				Id:         child.Id,
				Name:       child.Name,
				Properties: child.Properties,
				Type:       model.DeviceAction.Device.Type,
				Device:     &m.Device{Id: model.DeviceAction.Device.Id},
			}

			//*device = *model.DeviceAction.Device
			//device.Id = child.Id
			//device.Name = child.Name
			//device.Address = new(int)
			//*device.Address = *child.Address
			//device.Device = &m.Device{Id: model.DeviceAction.Device.Id}
			//device.Tty = child.Tty
			//device.Sleep = model.DeviceAction.Device.Sleep
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	nodes := f.core.GetNodes()
	nodeId := model.DeviceAction.Device.Node.Id
	if _, ok := nodes[nodeId]; ok {
		f.Node = nodes[nodeId]
	} else {
		// autoload nodes
		var node *m.Node
		if node, err = f.adaptors.Node.GetById(nodeId); err == nil {
			f.Node, _ = f.core.AddNode(node)
		} else {
			log.Error(err.Error())
			return
		}
	}

	// generate new worker
	worker := NewWorker(model, f, f.cron)

	// add devices to worker
	// ------------------------------------------------
	for _, device := range devices {

		var action *Action
		if action, err = NewAction(device, model.DeviceAction.Script, f.Node, f, f.scriptService); err != nil {
			log.Error(err.Error())
			continue
		}

		worker.AddAction(action)
	}

	f.Workers[model.Id] = worker
	f.Workers[model.Id].RegTask()

	return
}

func (f *Flow) UpdateWorker(worker *m.Worker) (err error) {

	if _, ok := f.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
	}

	if err = f.RemoveWorker(worker); err != nil {
		log.Warningf("error: %s", err.Error())
	}

	if err = f.AddWorker(worker); err != nil {
		log.Warningf("error: %s", err.Error())
	}

	return
}

func (f *Flow) RemoveWorker(worker *m.Worker) (err error) {

	log.Infof("Remove worker: \"%s\"", worker.Name)

	if _, ok := f.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
		return
	}

	// stop cron task
	f.Workers[worker.Id].RemoveTask()

	// delete worker
	delete(f.Workers, worker.Id)

	return
}

func (f *Flow) NewScript(s ...*m.Script) (engine *scripts.Engine, err error) {

	var model *m.Script
	if len(s) == 0 {
		model = &m.Script{
			Lang: ScriptLangJavascript,
		}
	} else {
		model = s[0]
	}

	if engine, err = f.workflow.NewScript(model); err != nil {
		return
	}

	javascript := engine.Get().(*scripts.Javascript)
	ctx := javascript.Ctx()

	// flow
	if b := ctx.GetGlobalString("IC"); !b {
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() *FlowBind {
		return &FlowBind{flow: f}
	})
	ctx.PutPropString(-3, "Flow")
	ctx.Pop()

	// workflow
	if b := ctx.GetGlobalString("IC"); !b {
		return
	}
	ctx.PushObject()
	ctx.PushGoFunction(func() *WorkflowBind {
		return &WorkflowBind{wf: f.workflow}
	})
	ctx.PutPropString(-3, "Workflow")
	ctx.Pop()

	return
}

func (f *Flow) defineCircularConnection(ctx context.Context) (newCtx context.Context, err error) {

	if v := ctx.Value("parents"); v != nil {
		if parents, ok := v.([]int64); ok {
			var exist bool
			for _, parentId := range parents {
				if parentId == f.Model.Id {
					exist = true
				}
			}

			if exist {
				depends := fmt.Sprintf("%d", parents[0])
				for _, parentId := range parents[1:] {
					depends = fmt.Sprintf("%s -> %d", depends, parentId)
				}
				err = fmt.Errorf("circular relationship detected: %s -> flow(%d)", depends, f.Model.Id)
				return
			}

			parents = append(parents, f.Model.Id)
			newCtx = context.WithValue(ctx, "parents", parents)

		} else {
			err = fmt.Errorf("bad parent context value: parents(%v)", parents)
		}

		return
	}

	newCtx = context.WithValue(ctx, "parents", []int64{f.Model.Id})

	return
}
