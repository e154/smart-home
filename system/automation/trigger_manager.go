package automation

import (
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/events"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/plugins/triggers"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"sync"
)

type triggerManager struct {
	eventBus       bus.Bus
	scriptService  scripts.ScriptService
	supervisor     supervisor.Supervisor
	adaptors       *adaptors.Adaptors
	isStarted      *atomic.Bool
	rawPlugin      triggers.IGetTrigger
	triggerCounter atomic.Uint64
	sync.Mutex
	triggers map[int64]*Trigger
}

func NewTriggerManager(eventBus bus.Bus,
	scriptService scripts.ScriptService,
	sup supervisor.Supervisor,
	adaptors *adaptors.Adaptors) (manager *triggerManager) {
	manager = &triggerManager{
		eventBus:      eventBus,
		scriptService: scriptService,
		supervisor:    sup,
		adaptors:      adaptors,
		isStarted:     atomic.NewBool(false),
		triggers:      make(map[int64]*Trigger),
	}
	return
}

// Start ...
func (a *triggerManager) Start() {

	a.load()
	_ = a.eventBus.Subscribe("system/automation/triggers/+", a.eventHandler)
	a.isStarted.Store(true)

	log.Info("Started")
}

// Shutdown ...
func (a *triggerManager) Shutdown() {

	a.unload()
	_ = a.eventBus.Unsubscribe("system/automation/triggers/+", a.eventHandler)

	log.Info("Shutdown")
}

func (a *triggerManager) eventHandler(_ string, msg interface{}) {

	switch v := msg.(type) {
	case events.EventEnableTrigger:
		go a.updateTrigger(v.Id)
	case events.EventDisableTrigger:
		go a.removeTrigger(v.Id)

	case events.EventUpdatedTrigger:
		go a.updateTrigger(v.Id)
	case events.EventAddedTrigger:
		go a.updateTrigger(v.Id)
	case events.EventRemovedTrigger:
		go a.removeTrigger(v.Id)
	}
}

func (a *triggerManager) load() {
	if a.isStarted.Load() {
		return
	}

	// load triggers plugin
	plugin, err := a.supervisor.GetPlugin(triggers.Name)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if rawPlugin, ok := plugin.(triggers.IGetTrigger); ok {
		a.rawPlugin = rawPlugin
	} else {
		log.Fatal("bad static cast triggers.IGetTrigger")
	}

	const perPage int64 = 500
	var page int64 = 0
LOOP:
	triggers, _, err := a.adaptors.Trigger.List(perPage, page*perPage, "", "", true)
	if err != nil {
		log.Error(err.Error())
		return
	}
	for _, trigger := range triggers {
		if err = a.addTrigger(trigger); err != nil {
			log.Warn(err.Error())
		}
	}
	if len(triggers) != 0 {
		page++
		goto LOOP
	}

	log.Info("Loaded ...")
}

func (a *triggerManager) unload() {
	if !a.isStarted.Load() {
		return
	}

	for id := range a.triggers {
		a.removeTrigger(id)
	}
	a.isStarted.Store(false)

	log.Info("Unloaded ...")
}

// addTrigger ...
func (a *triggerManager) addTrigger(model *m.Trigger) (err error) {

	defer func() {
		if err == nil {
			a.triggerCounter.Inc()
		}
	}()

	if _, ok := a.triggers[model.Id]; ok {
		err = errors.Wrap(apperr.ErrInternal, fmt.Sprintf("trigger %s exist", model.Name))
		return
	}

	var trigger *Trigger
	if trigger, err = NewTrigger(a.eventBus, a.scriptService, model, a.rawPlugin); err != nil {
		log.Error(err.Error())
		return
	}

	a.Lock()
	a.triggers[model.Id] = trigger
	a.Unlock()

	trigger.Start()

	return
}

// removeTrigger ...
func (a *triggerManager) removeTrigger(id int64) (err error) {
	a.Lock()
	defer a.Unlock()
	//log.Infof("remove trigger id:%d", id)

	trigger, ok := a.triggers[id]
	if !ok {
		return
	}
	trigger.Stop()
	delete(a.triggers, id)

	a.triggerCounter.Dec()
	return
}

// updateTrigger ...
func (a *triggerManager) updateTrigger(id int64) {
	//log.Infof("reload trigger id:%d", id)
	a.removeTrigger(id)

	task, err := a.adaptors.Trigger.GetById(id)
	if err != nil {
		return
	}

	a.addTrigger(task)
}

func (a *triggerManager) IsLoaded(id int64) (loaded bool) {
	a.Lock()
	_, loaded = a.triggers[id]
	a.Unlock()
	return
}
