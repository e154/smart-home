// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

package main

import (
	"fmt"
	"github.com/e154/smart-home/api/gate"
	"github.com/e154/smart-home/api/mobile"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/websocket"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"github.com/e154/smart-home/version"
	"os"
)

var (
	log = common.MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(version.ShortVersionBanner, version.GetHumanVersion())
			return
		case "-backup":
			container := BuildContainer()
			container.Invoke(func(
				backup *backup.Backup,
				graceful *graceful_service.GracefulService) {

				if err := backup.New(); err != nil {
					log.Error(err.Error())
				}

				graceful.Shutdown()
			})
			return
		case "-restore":
			if len(os.Args) < 3 {
				log.Error("need backup name")
				return
			}
			container := BuildContainer()
			container.Invoke(func(
				backup *backup.Backup,
				graceful *graceful_service.GracefulService) {

				if err := backup.Restore(os.Args[2]); err != nil {
					log.Error(err.Error())
				}

				graceful.Shutdown()
			})
			return
		case "-reset":
			container := BuildContainer()
			container.Invoke(func(
				initialService *initial.InitialService) {

				initialService.Reset()
			})
			return
		case "-demo":
			container := BuildContainer()
			container.Invoke(func(
				initialService *initial.InitialService) {

				initialService.InstallDemoData()
			})
			return
		default:
			fmt.Printf(version.VerboseVersionBanner, "v2", os.Args[0])
			return
		}
	}

	start()
}

func start() {

	fmt.Printf(version.ShortVersionBanner, "")

	container := BuildContainer()
	err := container.Invoke(func(server *server.Server,
		graceful *graceful_service.GracefulService,
		initialService *initial.InitialService,
		wsApi *websocket.WebSocket,
		mobileServer *mobile.MobileServer,
		metric *metrics.MetricManager,
		zigbee2mqtt *zigbee2mqtt.Zigbee2mqtt,
		gateApi *gate.Gate,
		logger *logging.Logging) {

		go server.Start()
		go mobileServer.Start()
		go wsApi.Start()
		go gateApi.Start()
		go metric.Start()
		go zigbee2mqtt.Start()

		graceful.Wait()
	})

	if err != nil {
		panic(err.Error())
	}
}
