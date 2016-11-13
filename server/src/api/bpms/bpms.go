package bpms

import (
	"log"
	"../models"
	cron "../crontab"
)

var (
	Cron        *cron.Crontab
)

type BPMS struct {
	nodes		[]*models.Node
	wfs		[]*Workflow
}

// Singleton
var instantiated *BPMS = nil

func BpmsPtr() *BPMS {
	return instantiated
}

func (b *BPMS) Init() (err error) {

	return
}

func (b *BPMS) Run() (err error) {

	if b.nodes, err = models.GetAllEnabledNodes(); err != nil {
		return
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

		b.wfs = append(b.wfs, wf)
	}

	for _, wf := range b.wfs {
		if err = wf.Run(); err != nil {
			return
		}
	}
	return
}

func (b *BPMS) Stop() (err error) {

	for _, wf := range b.wfs {
		if err = wf.Stop(); err != nil {
			return
		}
	}
	return
}

func (b *BPMS) Restart() (err error) {

	for _, wf := range b.wfs {
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