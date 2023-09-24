package local_migrations

import (
	"context"
	"strings"

	"github.com/e154/smart-home/adaptors"
	. "github.com/e154/smart-home/system/initial/assertions"
	"github.com/e154/smart-home/system/scripts"
)

type MigrationJavascript struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationJavascript(adaptors *adaptors.Adaptors) *MigrationJavascript {
	return &MigrationJavascript{
		adaptors: adaptors,
	}
}

func (n *MigrationJavascript) Up(ctx context.Context, adaptors *adaptors.Adaptors) error {
	if adaptors != nil {
		n.adaptors = adaptors
	}

	list, _, err := n.adaptors.Script.List(ctx, 999, 0, "desc", "id", nil)
	So(err, ShouldBeNil)

	var engine *scripts.Engine
	for _, script := range list {
		script.Source = strings.ReplaceAll(script.Source, "entityManager.getEntity", "GetEntity")

		script.Source = strings.ReplaceAll(script.Source, "entityManager.setState", "SetState ENTITY_ID,")
		script.Source = strings.ReplaceAll(script.Source, "supervisor.setState", "SetState ENTITY_ID,")
		script.Source = strings.ReplaceAll(script.Source, "entity.setState", "SetState ENTITY_ID,")
		script.Source = strings.ReplaceAll(script.Source, "Actor.setState", "SetState ENTITY_ID,")

		//script.Source = strings.ReplaceAll(script.Source, "supervisor.getEntity", "GetEntity")
		script.Source = strings.ReplaceAll(script.Source, "entityManager.setAttributes", "SetAttributes")
		script.Source = strings.ReplaceAll(script.Source, "entity.setAttributes", "SetAttributes")
		script.Source = strings.ReplaceAll(script.Source, "entity.getAttributes", "GetAttributes")
		script.Source = strings.ReplaceAll(script.Source, "entity.getAttributes()", "GetAttributes(ENTITY_ID)")
		script.Source = strings.ReplaceAll(script.Source, "GetAttributes()", "GetAttributes(ENTITY_ID)")
		script.Source = strings.ReplaceAll(script.Source, "entity.getSettings", "GetSettings")
		script.Source = strings.ReplaceAll(script.Source, "Actor.getSettings()", "GetSettings(ENTITY_ID)")
		script.Source = strings.ReplaceAll(script.Source, "GetSettings()", "GetSettings(ENTITY_ID)")

		script.Source = strings.ReplaceAll(script.Source, "entityManager.setMetric", "SetMetric")
		script.Source = strings.ReplaceAll(script.Source, "entityManager.callAction", "CallAction")
		script.Source = strings.ReplaceAll(script.Source, "entityManager.callScene", "CallScene")
		script.Source = strings.ReplaceAll(script.Source, "supervisor.setAttributes", "SetAttributes")

		script.Source = strings.ReplaceAll(script.Source, "supervisor.setMetric", "SetMetric")
		script.Source = strings.ReplaceAll(script.Source, "entity.setMetric", "SetMetric")

		script.Source = strings.ReplaceAll(script.Source, "supervisor.callAction", "CallAction")
		script.Source = strings.ReplaceAll(script.Source, "entity.callAction", "CallAction")
		script.Source = strings.ReplaceAll(script.Source, "Action.callAction", "CallAction")

		script.Source = strings.ReplaceAll(script.Source, "supervisor.callScene", "CallScene")
		engine, err = scripts.NewEngine(script, nil, nil)
		So(err, ShouldBeNil)

		err = engine.Compile()
		So(err, ShouldBeNil)

		err = n.adaptors.Script.Update(ctx, script)
		So(err, ShouldBeNil)
	}

	return nil
}
