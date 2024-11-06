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

//go:build !production

package commands

import (
	"fmt"

	"github.com/e154/smart-home/internal/common/app"
	"github.com/e154/smart-home/internal/system/initial"

	"github.com/spf13/cobra"
	"go.uber.org/fx"

	. "github.com/e154/smart-home/cmd/server/container"
	"github.com/e154/smart-home/version"
)

var (
	// Server ...
	Server = &cobra.Command{
		Use:   "server",
		Short: "Run smart home server",
		Run: func(cmd *cobra.Command, args []string) {

			mode = "test"
			fmt.Printf(version.ShortVersionBanner, version.VersionString, mode)

			app.Do(BuildContainer, fx.Invoke(func(
				_ *initial.Initial,
			) {

			}))

		},
	}
)

func init() {
	Server.AddCommand(backupCmd)
	Server.AddCommand(demoCmd)
	Server.AddCommand(resetCmd)
	Server.AddCommand(restoreCmd)
	Server.AddCommand(versionCmd)
	Server.AddCommand(gateCmd)
	Server.AddCommand(generateCertCmd)
}
