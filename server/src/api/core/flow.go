package core

import (
	"../log"
	"fmt"
	"errors"
	"github.com/astaxie/beego/orm"
	"../models"
	"encoding/hex"
)

func NewFlow(model *models.Flow, workflow *Workflow) (flow *Flow, err error) {

	flow = &Flow{
		Model: model,
		workflow: workflow,
		cursor: []*FlowElement{},
		quit: make(chan bool),
		Workers: make(map[int64]*Worker),
	}

	// get flow elements
	var flowelements []*models.FlowElement
	if flowelements, err = models.GetFlowElementsByFlow(model); err != nil {
		return
	}

	for _, element := range flowelements {
		flowElement, err := NewFlowElement(element,flow, workflow)
		if err == nil {
			flow.FlowElements = append(flow.FlowElements, flowElement)
		} else {
			log.Warn(err.Error())
		}
	}

	// get connections
	if flow.Connections, err = models.GetConnectionsByFlow(model); err != nil {
		return
	}

	for _, element := range flow.FlowElements {
		element.Flow = flow
		switch element.Model.PrototypeType  {
		case "MessageHandler":
			element.Prototype = &MessageHandler{}
			break
		case "MessageEmitter":
			element.Prototype = &MessageEmitter{}
			break
		case "Task":
			element.Prototype = &Task{}
			break
		case "Gateway":
			element.Prototype = &Gateway{}
			break
		case "Flow":
			element.Prototype = &FlowLink{}
			break
		}
	}

	go flow.loop()

	// add worker
	err = flow.InitWorkers()

	return
}

type Flow struct {
	Model        	*models.Flow
	workflow     	*Workflow
	Connections  	[]*models.Connection
	FlowElements 	[]*FlowElement
	Node		*models.Node
	Workers     	map[int64]*Worker
	cursor       	[]*FlowElement
	quit         	chan bool
}

func (f *Flow) Remove() {
	f.quit <- true
	for _, worker := range f.Workers {
		f.RemoveWorker(worker.Model)
	}
}

func (f *Flow) NewMessage(message *Message) (err error) {

	var exist bool
	for _, element := range f.FlowElements {
		if element.Prototype == nil {
			continue
		}

		if element.Model.PrototypeType != "MessageHandler" {
			continue
		}

		element.Run(message)

		exist = true
	}

	if !exist {
		err = errors.New("Message handler not found")
	}

	return
}

func (f *Flow) loop() {
	for {
		select {
		case <- f.quit:
			break
		}
	}
}

// ------------------------------------------------
// Workers
// ------------------------------------------------

func (f *Flow) InitWorkers() (err error) {

	var workers	[]*models.Worker
	if workers, err = f.Model.GetAllEnabledWorkers(); err != nil {
		return
	}

	for _, worker := range workers {
		if err = f.AddWorker(worker); err != nil {
			log.Warn(err.Error())
			return
		}
	}

	return
}

func (f *Flow) AddWorker(model *models.Worker) (err error) {

	log.Infof("Add worker: \"%s\"", model.Name)

	if _, ok := f.Workers[model.Id]; ok {
		return
	}

	if len(f.FlowElements) == 0 {
		err = errors.New("No flow elements")
		return
	}

	// generate new worker
	worker := NewWorker(model)

	// get device
	// ------------------------------------------------
	var devices []*models.Device
	if model.DeviceAction.Device.Address != nil {
		devices = append(devices, model.DeviceAction.Device)
	} else {
		// значит тут группа устройств
		var childs []*models.Device
		if childs, _, err = model.DeviceAction.Device.GetChilds(); err != nil {
			return
		}

		for _, child := range childs {
			if child.Address == nil || child.Status != "enabled" {
				continue
			}

			device := &models.Device{}
			*device = *model.DeviceAction.Device
			device.Id = child.Id
			device.Name = child.Name
			device.Address = new(int)
			*device.Address = *child.Address
			device.Device = &models.Device{Id:model.DeviceAction.Device.Id}
			device.Tty = child.Tty
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	if _, ok := f.workflow.Nodes[model.DeviceAction.Device.Node.Id]; ok {
		f.Node = f.workflow.Nodes[model.DeviceAction.Device.Node.Id]
	} else {
		// autoload nodes
		f.Node, err = models.GetNodeById(model.DeviceAction.Device.Node.Id)
		if err != nil {
			return
		}

		CorePtr().AddNode(f.Node)
	}

	// get script
	// ------------------------------------------------
	o := orm.NewOrm()
	if _, err = o.LoadRelated(model.DeviceAction, "Script"); err != nil {
		return
	}

	// add devices to worker
	// ------------------------------------------------
	for _, device := range devices {

		func(device *models.Device){

			var action *Action
			if action, err = NewAction(device, model.DeviceAction, f.Node); err != nil {
				return
			}

			action.Script.PushFunction("flow_new_message", func(v []byte) string {
				message := &Message{
					Result: hex.EncodeToString(v),
					Flow: f.Model,
					Device: device,
					Node: f.Node,
				}

				if err = f.NewMessage(message); err != nil {
					log.Warn(err.Error())
					return err.Error()
				}

				return ""
			})

			worker.AddAction(action)

		}(device)
	}

	f.Workers[model.Id] = worker
	f.Workers[model.Id].RegTask()

	return
}

func (f *Flow) UpdateWorker(worker *models.Worker) (err error) {

	if _, ok := f.Workers[worker.Id]; !ok {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
	}

	if err = f.RemoveWorker(worker); err != nil {
		log.Warn("error:", err.Error())
	}

	if err = f.AddWorker(worker); err != nil {
		log.Warn("error:", err.Error())
	}

	return
}

func (f *Flow) RemoveWorker(worker *models.Worker) (err error) {

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