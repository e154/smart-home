// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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

//
//         AUTOMATION
// ┌──────┐  ┌──────┐  ┌──────┐
// │ TASK │  │ TASK │  │ TASK │
// └──────┘  └──────┘  └──────┘   ...
//

package automation

import (
	"context"

	"github.com/e154/smart-home/internal/system/validation"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/common/telemetry"
	"github.com/e154/smart-home/pkg/events"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/plugins"
	"github.com/e154/smart-home/pkg/scripts"
	"go.uber.org/fx"

	"github.com/e154/bus"
)

var (
	log = logger.MustGetLogger("automation")
)

// Automation ...
type Automation interface {
	Start() (err error)
	Shutdown() (err error)
	Restart()
	TaskIsLoaded(id int64) bool
	TriggerIsLoaded(id int64) bool
	TaskTelemetry(id int64) telemetry.Telemetry
}

type automation struct {
	eventBus       bus.Bus
	triggerManager *triggerManager
	taskManager    *taskManager
}

// NewAutomation ...
func NewAutomation(lc fx.Lifecycle,
	eventBus bus.Bus,
	scriptService scripts.ScriptService,
	sup plugins.Supervisor,
	adaptors *adaptors.Adaptors,
	validation *validation.Validate) (auto Automation) {

	auto = &automation{
		eventBus:       eventBus,
		taskManager:    NewTaskManager(eventBus, scriptService, sup, adaptors),
		triggerManager: NewTriggerManager(eventBus, scriptService, sup, adaptors, validation),
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return auto.Shutdown()
		},
	})
	return
}

// Start ...
func (a *automation) Start() (err error) {
	a.taskManager.Start()
	a.triggerManager.Start()
	a.eventBus.Publish("system/services/automation", events.EventServiceStarted{Service: "Automation"})
	return
}

// Shutdown ...
func (a *automation) Shutdown() (err error) {
	a.taskManager.Shutdown()
	a.triggerManager.Shutdown()
	a.eventBus.Publish("system/services/automation", events.EventServiceStopped{Service: "Automation"})
	return
}

// Restart ...
func (a *automation) Restart() {
	_ = a.Shutdown()
	_ = a.Start()
}

func (a *automation) TaskIsLoaded(id int64) bool {
	return a.taskManager.IsLoaded(id)
}

func (a *automation) TriggerIsLoaded(id int64) bool {
	return a.triggerManager.IsLoaded(id)
}

func (a *automation) TaskTelemetry(id int64) telemetry.Telemetry {
	return a.taskManager.TaskTelemetry(id)
}
