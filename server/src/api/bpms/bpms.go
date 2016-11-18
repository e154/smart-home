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
	nodes     map[int64]*models.Node
	chanals	  map[int64]chan string
	workflows map[int64]*Workflow
}

func (b *BPMS) Init() (err error) {

	return
}

func (b *BPMS) Run() (err error) {

	if err = b.InitNodes(); err != nil {
		return
	}

	err = b.InitWorkflows()

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

func Initialize() (err error) {
	log.Println("BPMS initialize...")

	Cron = cron.CrontabPtr()

	bpmsPtr = &BPMS{}
	if err = bpmsPtr.Init(); err != nil {
		return
	}

	if err = bpmsPtr.Run(); err != nil {
		return
	}

	Hub = stream.GetHub()
	Hub.Subscribe("get.nodes.status", streamNodesStatus)

	return
}