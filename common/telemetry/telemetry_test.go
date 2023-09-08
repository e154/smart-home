package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/e154/smart-home/common/debug"
)

func TestTelemetry(t *testing.T) {

	ctx := context.Background()

	triggerCtx, triggerSpan := Start(ctx, "trigger")
	triggerSpan.SetAttributes("trigger.id", 1)
	time.Sleep(time.Millisecond * 500)
	triggerSpan.End()

	taskCtx, taskSpan := Start(triggerCtx, "task")
	taskSpan.SetAttributes("task.id", 1)

	conditionsCtx, conditionsSpan := Start(taskCtx, "conditions")
	conditionsSpan.SetAttributes("condition.id", 1)
	time.Sleep(time.Millisecond * 500)
	conditionsSpan.End()

	actionsCtx, actionsSpan := Start(conditionsCtx, "actions")
	actionsSpan.SetAttributes("action.id", 1)
	time.Sleep(time.Millisecond * 500)
	actionsSpan.End()

	taskSpan.End()

	//layers := Unpack(actionsCtx)
	//require.Equal(t, 2, len(layers))

	debug.Println(Unpack(actionsCtx))
}
