// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

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
