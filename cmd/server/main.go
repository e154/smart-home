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
	. "github.com/e154/smart-home/cmd/server/container"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/version"
	"go.uber.org/fx"
	"os"
)

var (
	log = MustGetLogger("main")
)

func main() {

	args := os.Args[1:]
	for _, arg := range args {
		switch arg {
		case "-v", "--version":
			fmt.Printf(version.ShortVersionBanner, version.GetHumanVersion())
			return
		case "-backup":
			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				backup *backup.Backup) {

				if err := backup.New(); err != nil {
					log.Error(err.Error())
				}
			}))
			Start(app)
			return
		case "-restore":
			if len(os.Args) < 3 {
				log.Error("need backup name")
				return
			}
			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				backup *backup.Backup) {

				if err := backup.Restore(os.Args[2]); err != nil {
					log.Error(err.Error())
				}

			}))
			Start(app)
			return
		case "-reset":
			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				initialService *initial.Initial) {

				initialService.Reset()
			}))
			Start(app)
			return
		case "-demo":
			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				initialService *initial.Initial) {

				initialService.InstallDemoData()
			}))
			Start(app)
			return
		default:
			fmt.Printf(version.VerboseVersionBanner, "v2", os.Args[0])
			return
		}
	}

	app := BuildContainer(fx.Invoke(func(
		logger *logging.Logging,
		dbSaver logging.ISaver,
		initialService *initial.Initial,
	) {
		logger.SetSaver(dbSaver)
	}))

	Start(app)

	Work()

	Stop(app)
}
