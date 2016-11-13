package bpms

import (
	"log"
	"../models"
	cron "../crontab"
)

var (
	Cron        *cron.Crontab
)

// Singleton
var instantiated *BPMS = nil

func BpmsPtr() *BPMS {
	return instantiated
}

type BPMS struct {
	nodes     map[int64]*models.Node
	workflows map[int64]*Workflow
}

func (b *BPMS) Init() (err error) {

	return
}

func (b *BPMS) Run() (err error) {

	var nodes []*models.Node
	b.nodes = make(map[int64]*models.Node)
	b.workflows = make(map[int64]*Workflow)

	if nodes, err = models.GetAllEnabledNodes(); err != nil {
		return
	}
	for _, node := range nodes {
		b.nodes[node.Id] = node
	}

	log.Println("--------------------- NODES ---------------------")
	for _, node := range b.nodes {
		if _, err := node.RpcDial(); err != nil {
			log.Printf("Node error %s", err.Error())
			continue
		}

		log.Printf("Node dial tcp %s:%d ... ok",node.Ip, node.Port)
	}

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

func (b *BPMS) AddNode(node *models.Node) *models.Node {
	if _, ok := b.nodes[node.Id]; !ok {
		b.nodes[node.Id] = node
	}
	
	return b.nodes[node.Id]
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

func Initialize() (err error) {
	log.Println("BPMS initialize...")

	Cron = cron.CrontabPtr()

	instantiated = &BPMS{}
	if err = instantiated.Init(); err != nil {
		return
	}

	if err = instantiated.Run(); err != nil {
		return
	}

	return
}