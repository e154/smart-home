package main

import (
	"./workflow"
	"log"
	"time"
)

func main() {

	flow := &workflow.Flow{
		Name: "Какой-то бизнес процесс",
		Type: "generalized",
		Created_at: time.Now(),
	}

	flow_element := &workflow.FlowElement{
		Flow: flow,
	}

	log.Println(flow_element)
}