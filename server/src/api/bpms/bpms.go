package bpms

import (
	"log"
	"../models"
	"../stream"
	cron "../crontab"
)

var (
	Cron		*cron.Crontab
	Hub		stream.Hub
	bpmsPtr         *BPMS = nil
)

func BpmsPtr() *BPMS {
	return bpmsPtr
}

type BPMS struct {
	nodes      map[int64]*models.Node
	nodes_chan map[int64]chan string
	workflows  map[int64]*Workflow
}

func (b *BPMS) Run() (err error) {

	if err = b.InitNodes(); err != nil {
		return
	}

	err = b.InitWorkflows()

	return
}

func (b *BPMS) InitNodes() (err error) {
	var nodes []*models.Node
	b.nodes = make(map[int64]*models.Node)
	b.nodes_chan = make(map[int64]chan string)

	//log.Println("--------------------- NODES ---------------------")
	if nodes, err = models.GetAllEnabledNodes(); err != nil {
		return
	}

	for _, node := range nodes {
		b.AddNode(node)
	}

	//TODO remove
	if len(b.nodes) == 0 {
		return
	}

	return
}

func (b *BPMS) InitWorkflows() (err error) {

	b.workflows = make(map[int64]*Workflow)
	//log.Println("------------------- WORKFLOW --------------------")
	workflows, err := models.GetAllEnabledWorkflow()
	if err != nil {
		return
	}

	for _, workflow := range workflows {
		b.AddWorkflow(workflow)
	}

	return
}

func (b *BPMS) Stop() (err error) {

	for _, wf := range b.workflows {
		if err = wf.Stop(); err != nil {
			return
		}
	}
	return
}

func (b *BPMS) Restart() (err error) {

	for _, wf := range b.workflows {
		if err = wf.Restart(); err != nil {
			return
		}
	}
	return
}

func (b *BPMS) AddFlow(f *models.Flow) (err error) {

	var flow *models.Flow
	if flow, err = models.GetFlowById(f.Id); err != nil {
		return
	}

	if _, ok := b.workflows[flow.WorkflowId]; ok {
		if err = b.workflows[flow.WorkflowId].AddFlow(flow); err != nil {
			return
		}
	}

	return
}

func (b *BPMS) UpdateFlow(f *models.Flow) (err error) {

	var flow *models.Flow
	if flow, err = models.GetFlowById(f.Id); err != nil {
		return
	}

	if _, ok := b.workflows[flow.WorkflowId]; ok {
		if err = b.workflows[flow.WorkflowId].UpdateFlow(flow); err != nil {
			return
		}
	}

	return
}

func (b *BPMS) RemoveFlow(f *models.Flow) (err error) {

	var flow *models.Flow
	if flow, err = models.GetFlowById(f.Id); err != nil {
		return
	}

	if _, ok := b.workflows[flow.WorkflowId]; ok {
		if err = b.workflows[flow.WorkflowId].RemoveFlow(flow); err != nil {
			return
		}
	}

	return
}

func Initialize() (err error) {
	log.Println("BPMS initialize...")

	Cron = cron.CrontabPtr()

	bpmsPtr = &BPMS{}
	if err = bpmsPtr.Run(); err != nil {
		return
	}

	Hub = stream.GetHub()
	Hub.Subscribe("get.nodes.status", streamNodesStatus)
	Hub.Subscribe("get.workflow.status", streamWorkflowsStatus)
	Hub.Subscribe("get.flows.status", streamWorkflowsStatus)

	return
}