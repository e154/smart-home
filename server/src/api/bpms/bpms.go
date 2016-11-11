package bpms

import (
	"log"
	"../workflow"
)


type BPMS struct {
	wfs		[]*workflow.Workflow
}

// Singleton
var instantiated *BPMS = nil

func BpmsPtr() *BPMS {
	return instantiated
}

func (b *BPMS) Init() (err error) {

	b.wfs, err = workflow.Initialize()

	return
}

func (b *BPMS) Run() (err error) {

	return
}

func (b *BPMS) Stop() (err error) {

	return
}

func (b *BPMS) Restart() (err error) {

	return
}

func Initialize() error {
	log.Println("BPMS initialize...")

	instantiated = &BPMS{}
	return instantiated.Init()
}