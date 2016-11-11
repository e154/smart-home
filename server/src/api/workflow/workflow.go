package workflow

import (
	"../models"
	"encoding/json"
	"log"
)

type Workflow struct {
	model		*models.Workflow
	flows	 	[]*models.Flow
}

func (wf *Workflow) Init() (err error) {

	wf.flows, err = wf.model.GetAllEnabledFlows()
	if err != nil {
		return
	}

	// debug
	for _, flow := range wf.flows {
		j, _ := json.Marshal(flow)
		log.Println("--------------------- SQL ---------------------")
		log.Println(string(j))
		log.Println("--------------------- SQL ---------------------")
		flow.NewMessage(&models.Message{})
	}

	return
}

func Initialize() (wfs_constrollers []*Workflow, err error) {

	wfs_models, err := models.GetAllEnabledWorkflow()
	if err != nil {
		return
	}

	for _, wf_model := range wfs_models {

		wf := &Workflow{model: wf_model}
		if err = wf.Init(); err != nil {
			return
		}

		wfs_constrollers = append(wfs_constrollers, wf)
	}

	return
}