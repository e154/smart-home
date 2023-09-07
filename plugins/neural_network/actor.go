package neural_network

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
)

type Actor struct {
	supervisor.BaseActor
	eventBus      bus.Bus
	adaptors      *adaptors.Adaptors
	scriptService scripts.ScriptService
	actionPool    chan events.EventCallEntityAction
	network1      *Network1
	network2      *Network2
}

func NewActor(entity *m.Entity,
	visor supervisor.Supervisor,
	adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService,
	eventBus bus.Bus) *Actor {

	actor := &Actor{
		BaseActor:     supervisor.NewBaseActor(entity, scriptService, adaptors),
		adaptors:      adaptors,
		scriptService: scriptService,
		eventBus:      eventBus,
		actionPool:    make(chan events.EventCallEntityAction, 10),
		network1:      NewNetwork1(eventBus),
		network2:      NewNetwork2(eventBus, visor),
	}

	actor.Supervisor = visor

	// Actions
	for _, a := range actor.Actions {
		if a.ScriptEngine != nil {
			// bind
			a.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
			_, _ = a.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
			_, _ = a.ScriptEngine.Do()
		}
	}

	if actor.ScriptEngine != nil {
		_, _ = actor.ScriptEngine.EvalString(fmt.Sprintf("const ENTITY_ID = \"%s\";", entity.Id))
		actor.ScriptEngine.PushStruct("Actor", supervisor.NewScriptBind(actor))
	}

	// action worker
	go func() {
		for msg := range actor.actionPool {
			actor.runAction(msg)
		}
	}()

	return actor
}

func (e *Actor) destroy() {
	if e.network2 != nil {
		e.network2.Stop()
	}
}

func (e *Actor) Spawn() supervisor.PluginActor {
	e.Update()
	return e
}

func (e *Actor) Update() {

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

func (e *Actor) Start() {
}

func (e *Actor) Stop() {
}
