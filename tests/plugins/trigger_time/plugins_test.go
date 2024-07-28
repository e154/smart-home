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

package trigger_time

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/e154/bus"
	"go.uber.org/dig"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/supervisor"
	. "github.com/e154/smart-home/tests/plugins"
	. "github.com/e154/smart-home/tests/plugins/container"
)

func init() {
	apppath := filepath.Join(os.Getenv("PWD"), "../../..")
	_ = os.Chdir(apppath)
}

var (
	container *dig.Container
)

func TestMain(m *testing.M) {

	_ = os.Setenv("TEST_MODE", "true")

	container = BuildContainer()
	err := container.Invoke(func(
		_ *logging.Logging,
		migrations *migrations.Migrations,
		adaptors *adaptors.Adaptors,
		scriptService scripts.ScriptService,
		supervisor supervisor.Supervisor,
		automation automation.Automation,
		scheduler *scheduler.Scheduler,
		eventBus bus.Bus,
	) {

		migrations.Purge()

		time.Sleep(time.Millisecond * 500)

		// register plugins
		AddPlugin(adaptors, "triggers")
		AddPlugin(adaptors, "time")
		AddPlugin(adaptors, "sensor")

		serviceCh := WaitService(eventBus, time.Second*5, "Scheduler", "Automation", "Supervisor")
		pluginsCh := WaitPlugins(eventBus, time.Second*5, "triggers", "sensor", "time")
		scriptService.Restart()
		scheduler.Start(context.Background())
		supervisor.Start(context.Background())
		automation.Start()

		defer supervisor.Shutdown(context.Background())
		defer automation.Shutdown()

		<-serviceCh
		<-pluginsCh

		os.Exit(m.Run())
	})

	if err != nil {
		fmt.Println("error:", dig.RootCause(err))
	}
}
