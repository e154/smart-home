package bpms

import (
	"../models"
	r "../../lib/rpc"
	"log"
	"time"
	"encoding/hex"
)

type Workflow struct {
	model		*models.Workflow
	flows	 	map[int64]*models.Flow
	workers		map[int64]*models.Worker
	nodes		map[int64]*models.Node
}

func (b *BPMS) InitWorkflows() (err error) {

	b.workflows = make(map[int64]*Workflow)
	log.Println("------------------- WORKFLOW --------------------")
	workflows, err := models.GetAllEnabledWorkflow()
	if err != nil {
		return
	}
	log.Println("ok")

	for _, workflow := range workflows {

		wf := &Workflow{model: workflow, nodes: b.nodes}
		if err = wf.Init(); err != nil {
			return
		}

		b.workflows[workflow.Id] = wf
	}

	for _, wf := range b.workflows {
		if err = wf.Run(); err != nil {
			return
		}
	}

	return
}

func (wf *Workflow) Init() (err error) {

	var flows	[]*models.Flow
	var workers	[]*models.Worker

	wf.flows = make(map[int64]*models.Flow)
	wf.workers = make(map[int64]*models.Worker)

	log.Println("-------------------- FLOWS ----------------------")
	if flows, err = wf.model.GetAllEnabledFlows(); err != nil {return}
	for _, flow := range flows {
		wf.flows[flow.Id] = flow
	}
	log.Println("ok")

	log.Println("------------------- WORKERS ---------------------")
	if workers, err = models.GetAllEnabledWorkers(); err != nil {return}
	for _, worker := range workers {
		wf.workers[worker.Id] = worker
	}

	for _, worker := range wf.workers {
		wf.AddWorker(worker)
	}
	log.Println("ok")

	return
}

func (wf *Workflow) Run() (err error) {

	for _, worker := range wf.workers {
		worker.CronTask.Run()
	}

	return
}

func (wf *Workflow) Stop() (err error) {

	for _, worker := range wf.workers {
		worker.CronTask.Stop()
	}

	return
}

func (wf *Workflow) Restart() (err error) {

	wf.Stop()
	wf.Run()

	return
}

func (wf *Workflow) AddWorker(worker *models.Worker) (err error) {
	log.Printf("start \"%s\"", worker.Name)
	//j, _ := json.Marshal(worker)
	//log.Println(string(j))

	// autoload flows
	if _, ok := wf.flows[worker.FlowId]; ok {
		worker.Flow = wf.flows[worker.FlowId]
	} else {
		var flow *models.Flow
		flow, err = models.GetEnabledFlowById(worker.FlowId)
		if err != nil {
			return
		}

		wf.flows[flow.Id] = flow
		worker.Flow = wf.flows[flow.Id]
	}

	// message
	worker.Message = &models.Message{Variable: []byte(worker.DeviceAction.Command)}

	// autoload nodes
	if _, ok := wf.nodes[worker.Device.NodeId]; ok {
		worker.Node = wf.nodes[worker.Device.NodeId]
	} else {
		var node *models.Node
		node, err = models.GetNodeById(worker.Device.NodeId)
		if err != nil {
			return
		}

		BpmsPtr().AddNode(node)
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

	worker.CronTask = Cron.NewTask(worker.Time, func() {

		args.Time = time.Now()

		result := &r.Result{}
		if !worker.Node.IsConnected() {
			worker.Node.Errors++
			return
		}

		if err := worker.Node.ModbusSend(args, result); err != nil {
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