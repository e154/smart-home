package example1

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/scripts"
)

type Example1 struct {
	adaptors       *adaptors.Adaptors
	scriptService  scripts.ScriptService
	areaManager    *AreaManager
	scriptManager  *ScriptManager
	supervisor     *Supervisor
	triggerManager *TriggerManager
}

func NewExample1(adaptors *adaptors.Adaptors,
	scriptService scripts.ScriptService) *Example1 {
	return &Example1{
		adaptors:       adaptors,
		scriptService:  scriptService,
		areaManager:    NewAreaManager(adaptors),
		scriptManager:  NewScriptManager(adaptors, scriptService),
		supervisor:     NewSupervisor(adaptors),
		triggerManager: NewTriggerManager(adaptors),
	}
}

func (e *Example1) Install(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		e.adaptors = adaptors
		e.areaManager = NewAreaManager(adaptors)
		e.scriptManager = NewScriptManager(adaptors, e.scriptService)
		e.supervisor = NewSupervisor(adaptors)
		e.triggerManager = NewTriggerManager(adaptors)
	}

	var area *m.Area
	if area, err = e.areaManager.addArea(ctx, "zone51", "demo"); err != nil {
		return
	}
	var scripts []*m.Script
	if scripts, err = e.scriptManager.addScripts(ctx); err != nil {
		return
	}

	var entities []*m.Entity
	if entities, err = e.supervisor.addEntities(ctx, scripts, area); err != nil {
		return
	}

	e.triggerManager.addTriggers(scripts, entities)

	return
}
