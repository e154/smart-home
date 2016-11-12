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
	flows	 	[]*models.Flow
	workers		[]*models.Worker
	nodes		[]*models.Node


}

func (wf *Workflow) Init() (err error) {

	var flows	[]*models.Flow
	var workers	[]*models.Worker

	log.Println("-------------------- FLOWS ----------------------")
	if flows, err = wf.model.GetAllEnabledFlows(); err != nil {
		return
	}
	log.Println("ok")

	log.Println("------------------- WORKERS ---------------------")
	if workers, err = models.GetAllEnabledWorkers(); err != nil {
		return
	}

	for _, worker := range workers {
		log.Printf("start \"%s\"", worker.Name)
		//j, _ := json.Marshal(worker)
		//log.Println(string(j))

		// flows
		for _, flow := range flows {
			if flow.Id == worker.FlowId {
				worker.Flow = flow
				break
			}
		}

		// message
		worker.Message = &models.Message{Variable: []byte(worker.DeviceAction.Command)}

		// flows
		for _, node := range wf.nodes {
			if node.Id == worker.Device.NodeId {
				worker.Node = node
			}
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

		WM.Run(worker.Time, func() {

			args.Time = time.Now()

			result := &r.Result{}
			if err := worker.Node.ModbusSend(args, result); err != nil {
				log.Println("err ", err.Error())
			}

			worker.Message.Variable = result.Result
			if err := worker.Flow.NewMessage(worker.Message); err != nil {
				log.Println("err" , err.Error())
			}
		})

	}
	log.Println("ok")

	wf.flows = flows

	return
}

func (wf *Workflow) Run() (err error) {

	return
}

func (wf *Workflow) Stop() (err error) {

	return
}

func (wf *Workflow) Restart() (err error) {

	return
}
