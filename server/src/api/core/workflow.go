package core

import (
	"log"
	"time"
	"reflect"
	"encoding/hex"
	"github.com/astaxie/beego/orm"
	cr "github.com/e154/cron"
	r "../../lib/rpc"
	"encoding/json"
	"errors"
	"fmt"
	"../models"
	"../stream"
	"../cron"
)

type Workflow struct {
	model   *models.Workflow
	Flows   map[int64]*models.Flow
	Workers map[int64]*models.Worker
	CronTasks map[int64]*cr.Task
	Nodes   map[int64]*models.Node
}

func (wf *Workflow) Run() (err error) {

	err = wf.InitFlows()
	if err != nil {
		return
	}

	err = wf.InitWorkers()

	return
}

func (wf *Workflow) InitFlows() (err error) {

	//log.Println("-------------------- FLOWS ----------------------")

	var flows	[]*models.Flow
	wf.Flows = make(map[int64]*models.Flow)
	if flows, err = wf.model.GetAllEnabledFlows(); err != nil {return}
	for _, flow := range flows {
		wf.AddFlow(flow)
	}

	return
}

func (wf *Workflow) InitWorkers() (err error) {

	//log.Println("------------------- WORKERS ---------------------")

	var workers	[]*models.Worker
	wf.Workers = make(map[int64]*models.Worker)
	if workers, err = models.GetAllEnabledWorkersByWorkflow(wf.model); err != nil {return}
	for _, worker := range workers {
		if err = wf.AddWorker(worker); err != nil {
			log.Println("error:", err.Error())
			return
		}

	}

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, task := range wf.CronTasks {
		task.Disable()
	}

	return
}

func (wf *Workflow) Restart() (err error) {

	wf.Stop()
	err = wf.Run()

	return
}

func (wf *Workflow) AddFlow(flow *models.Flow) (err error) {

	if flow.Status != "enabled" {
		return
	}

	log.Println("Add flow:", flow.Name)

	if _, ok := wf.Flows[flow.Id]; ok {
		return
	}

	wf.Flows[flow.Id] = flow

	return
}

func (wf *Workflow) UpdateFlow(flow *models.Flow) (err error) {

	if err = wf.RemoveFlow(flow); err != nil {
		return
	}

	return wf.AddFlow(flow)
}

func (wf *Workflow) RemoveFlow(flow *models.Flow) (err error) {

	//TODO remove workers for this flow

	log.Println("Remove flow:", flow.Name)

	if _, ok := wf.Flows[flow.Id]; !ok {
		return
	}

	delete(wf.Flows, flow.Id)

	return
}

func (wf *Workflow) AddWorker(worker *models.Worker) (err error) {

	log.Printf("Add worker: \"%s\"", worker.Name)

	j, _ := json.Marshal(worker)
	log.Println("worker:",string(j))

	if worker.DeviceAction != nil {
		o := orm.NewOrm()
		o.LoadRelated(worker, "DeviceAction")
	}

	if worker.DeviceAction.Device != nil {
		o := orm.NewOrm()
		o.LoadRelated(worker.DeviceAction, "Device")
	}

	if worker.Device, err = models.GetParentDeviceByChildId(worker.DeviceAction.Device.Id); err != nil {
		return
	}

	if worker.Device == nil {
		err = errors.New("device not found")
		return
	}

	o := orm.NewOrm()
	if _, err = o.LoadRelated(worker.Device, "Node"); err != nil {
		return
	}

	worker.Device.Id = worker.DeviceAction.Device.Id

	if _, ok := wf.Workers[worker.Id]; ok {
		return
	}

	wf.Workers[worker.Id] = worker

	//j, _ := json.Marshal(worker)
	//log.Println(string(j))

	// autoload flows
	if _, ok := wf.Flows[worker.Flow.Id]; ok {
		worker.Flow = wf.Flows[worker.Flow.Id]
	} else {
		var flow *models.Flow
		flow, err = models.GetEnabledFlowById(worker.Flow.Id)
		if err != nil {
			return
		}

		wf.Flows[flow.Id] = flow
		worker.Flow = wf.Flows[flow.Id]
	}

	// message
	worker.Message = &models.Message{Variable: []byte(worker.DeviceAction.Command)}

	if worker.Device.Node.Id == 0 {
		err = errors.New("device node.id = 0")
		return
	}

	// check node list
	if _, ok := wf.Nodes[worker.Device.Node.Id]; ok {
		worker.Node = wf.Nodes[worker.Device.Node.Id]
	} else {
		// autoload nodes
		var node *models.Node
		node, err = models.GetNodeById(worker.Device.Node.Id)
		if err != nil {
			return
		}

		CorePtr().AddNode(node)
		worker.Node = node
	}

	command, _ := hex.DecodeString(worker.DeviceAction.Command)

	args := r.Request{
		Baud: worker.Device.Baud,
		Result: true,
		Command: command,
		Device: worker.Device.Tty,
		Line: "",
		StopBits: 2,
		Time: time.Now(),
		Timeout: worker.Device.Timeout,
	}

	wf.CronTasks[worker.Id] = cron.Cron().NewTask(worker.Time, func() {

		args.Time = time.Now()

		result := &r.Result{}
		if !worker.Node.IsConnected() {
			worker.Node.Errors++
			return
		}

		if err := worker.Node.ModbusSend(args, result); err != nil {
			log.Println(err.Error())
			worker.Node.Errors++
			// нет связи с нодой, или что-то случилось
			return
		}

		worker.Message.Variable = result.Result
		if err := worker.Flow.NewMessage(worker.Message); err != nil {
			log.Println("err" , err.Error())
		}
	})

	return
}

func (wf *Workflow) UpdateWorker(worker *models.Worker) (err error) {

	log.Printf("Update worker: \"%s\"", worker.Name)

	if _, ok := wf.Workers[worker.Id]; ok {

		wf.RemoveWorker(worker)

		if err = wf.AddWorker(worker); err != nil {
			log.Println("error:", err.Error())
		}

	} else {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
	}

	return
}

func (wf *Workflow) RemoveWorker(worker *models.Worker) (err error) {

	log.Printf("Remove worker: \"%s\"", worker.Name)

	if _, ok := wf.Workers[worker.Id]; ok {

		// stop cron task
		wf.CronTasks[worker.Id].Disable()

		// remove task from cron
		cron.Cron().RemoveTask(wf.CronTasks[worker.Id])

		// delete worker
		delete(wf.Workers, worker.Id)
		delete(wf.CronTasks, worker.Id)

	} else {
		err = fmt.Errorf("worker id:%d not found", worker.Id)
	}

	return
}

func (wf *Workflow) GetStatus() string {
	return wf.model.Status
}

func GetWorkflowsStatus() (result map[int64]string) {
	result = make(map[int64]string)
	for id, workflow := range corePtr.workflows {
		result[id] = workflow.model.Status
	}

	return
}

func streamWorkflowsStatus(client *stream.Client, value interface{}) {
	v, ok := reflect.ValueOf(value).Interface().(map[string]interface{})
	if !ok {
		return
	}

	result := GetWorkflowsStatus()
	msg, _ := json.Marshal(map[string]interface{}{"id": v["id"], "workflows": result})
	client.Send(string(msg))
}