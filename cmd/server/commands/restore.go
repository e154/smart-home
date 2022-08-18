// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package commands

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	. "github.com/e154/smart-home/cmd/server/container"
	. "github.com/e154/smart-home/common/app"
	. "github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/logging"
)

var (
	log = MustGetLogger("main")
)

var (
	filename string

	restoreCmd = &cobra.Command{
		Use:   "restore",
		Short: "Restore settings from backup archive",
		Run: func(cmd *cobra.Command, args []string) {

			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				backup *backup.Backup) {

				if err := backup.Restore(filename); err != nil {
					log.Error(err.Error())
				}

			}))
			Start(app)
		},
	}
)

func init() {
	backupCmd.Flags().StringVarP(&filename, "filename", "f", "backup.zip", "backup file name")
}
