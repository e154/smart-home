package neural_network

import (
	"fmt"

	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

type Actor struct {
	supervisor.BaseActor
	actionPool chan events.EventCallEntityAction
	network1   *Network1
	network2   *Network2
}

func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	actor := &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		actionPool: make(chan events.EventCallEntityAction, 10),
		network1:   NewNetwork1(service.EventBus()),
		network2:   NewNetwork2(service.EventBus()),
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

func (e *Actor) Destroy() {
	if e.network2 != nil {
		e.network2.Stop()
	}
}

func (e *Actor) Spawn() {
}

func (e *Actor) addAction(event events.EventCallEntityAction) {
	e.actionPool <- event
}

func (e *Actor) runAction(msg events.EventCallEntityAction) {
	action, ok := e.Actions[msg.ActionName]
	if !ok {
		log.Warnf("action %s not found", msg.ActionName)
		return
	}
	switch action.Name {
	case "TRAIN1":
		e.network2.Train1()
	case "TRAIN2":
		e.network2.Train2()
	case "TRAIN3":
	case "TRAIN4":
	case "CHECK2":
	case "ENABLE":
	case "DISABLE":

	default:
		fmt.Sprintf("unknown comand: %s", action.Name)
	}
}
