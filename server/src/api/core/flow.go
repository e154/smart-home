package core

import (
	"../log"
	"fmt"
	"errors"
	"github.com/astaxie/beego/orm"
	"../scripts"
	cr "github.com/e154/cron"
	r "../../lib/rpc"
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

func (f *Flow) AddWorker(worker *models.Worker) (err error) {

	log.Infof("Add worker: \"%s\"", worker.Name)

	if _, ok := f.Workers[worker.Id]; ok {
		return
	}

	if len(f.FlowElements) == 0 {
		err = errors.New("No flow elements")
		return
	}

	f.Workers[worker.Id] = &Worker{Model:worker,}

	// get device
	// ------------------------------------------------
	var devices []*models.Device
	if worker.DeviceAction.Device.Address != nil {
		devices = append(devices, worker.DeviceAction.Device)
	} else {
		// значит тут группа устройств
		var childs []*models.Device
		if childs, _, err = worker.DeviceAction.Device.GetChilds(); err != nil {
			return
		}

		for _, child := range childs {
			if child.Address == nil || child.Status != "enabled" {
				continue
			}

			device := &models.Device{}
			*device = *worker.DeviceAction.Device
			device.Id = child.Id
			device.Name = child.Name
			device.Address = new(int)
			*device.Address = *child.Address
			device.Device = &models.Device{Id:worker.DeviceAction.Device.Id}
			device.Tty = child.Tty
			devices = append(devices, device)
		}
	}

	// get node
	// ------------------------------------------------
	if _, ok := f.workflow.Nodes[worker.DeviceAction.Device.Node.Id]; ok {
		f.Node = f.workflow.Nodes[worker.DeviceAction.Device.Node.Id]
	} else {
		// autoload nodes
		f.Node, err = models.GetNodeById(worker.DeviceAction.Device.Node.Id)
		if err != nil {
			return
		}

		CorePtr().AddNode(f.Node)
	}

	// cron worker
	// ------------------------------------------------
	for _, device := range devices {

		func(device *models.Device){

			// device
			if f.Workers[worker.Id].Devices == nil {
				f.Workers[worker.Id].Devices = make(map[int64]*models.Device)
			}

			f.Workers[worker.Id].Devices[device.Id] = device

			// get script
			// ------------------------------------------------
			o := orm.NewOrm()
			if _, err = o.LoadRelated(worker.DeviceAction, "Script"); err != nil {
				return
			}

			// add script
			script, _ := scripts.New(worker.DeviceAction.Script)
			script.PushStruct("device", device)
			script.PushStruct("flow", f.Model)
			script.PushStruct("node", f.Node)
			script.PushStruct("request", &r.Request{})

			script.PushFunction("modbus_send", func(args *r.Request) (result r.Result) {
				if err := f.Node.ModbusSend(args, &result); err != nil {
					result.Error = err.Error()
				}

				return
			})

			script.PushFunction("flow_new_message", func(v []byte) string {
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

			// set cron
			// ------------------------------------------------
			if f.Workers[worker.Id].CronTasks == nil {
				f.Workers[worker.Id].CronTasks = make(map[int64]*cr.Task)
			}

			f.Workers[worker.Id].CronTasks[device.Id] = cron.NewTask(worker.Time, func() {
				script.Do()
			})
		}(device)
	}

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
	for _, task := range f.Workers[worker.Id].CronTasks {

		task.Disable()

		// remove task from cron
		cron.RemoveTask(task)
	}

	// delete worker
	delete(f.Workers, worker.Id)

	return
}