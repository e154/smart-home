// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package terminal

import (
	"bytes"
	"context"
	"sync"

	"github.com/spf13/cobra"

	"github.com/e154/smart-home/common/events"
)

type CommandSudo struct {
	Out *bytes.Buffer
	*sync.Mutex
	cmd      *cobra.Command
	password *string
}

func NewCommandSudo() *CommandSudo {
	buf := new(bytes.Buffer)

	command := &CommandSudo{
		Out:      buf,
		password: new(string),
		Mutex:    &sync.Mutex{},
	}

	cmdMode := &cobra.Command{
		Use:               "sudo",
		DisableAutoGenTag: true,
		Short:             "get superuser access rights",
		Run: func(cmd *cobra.Command, args []string) {
			var password = cmd.Flag("password").Value.String()

			log.Debugf("password: %s", password)

			ctx := cmd.Context()
			command, ok := ctx.Value("Command").(events.CommandTerminal)
			if !ok {
				cmd.OutOrStdout().Write([]byte("request rejected"))
				return
			}

			options, ok := ctx.Value("Options").(CommandOptions)
			if !ok {
				cmd.OutOrStdout().Write([]byte("request rejected"))
				return
			}

			if options.config.RootSecret == "" || options.config.RootSecret != password {
				cmd.OutOrStdout().Write([]byte("request rejected"))
				return
			}

			accessToken, err := options.jwtManager.Generate(command.User, true)
			if err != nil {
				cmd.OutOrStdout().Write([]byte("request rejected"))
				log.Error(err.Error())
				return
			}

			cmd.OutOrStdout().Write([]byte("request processed successfully"))
			options.eventBus.Publish("system/dashboard", events.EventDirectMessage{
				UserID:    command.User.Id,
				SessionID: command.SessionID,
				Query:     "update_access_token",
				Message: UpdateAccessTokenMessage{
					AccessToken: accessToken,
				},
			})
		},
	}

	cmdMode.SetOut(command.Out)
	cmdMode.SetErr(command.Out)
	command.cmd = cmdMode

	return command
}

func (c *CommandSudo) Execute(ctx context.Context, args ...string) (out string, err error) {
	c.Lock()
	defer c.Unlock()
	c.Out.Reset()
	c.cmd.ResetFlags()
	c.cmd.Flags().StringVarP(c.password, "password", "p", "", "see the root_secret in configuration file or ROOT_SECRET in your environment")
	c.cmd.SetArgs(args)
	err = c.cmd.ExecuteContext(ctx)
	out = c.Out.String()
	return
}

func (c *CommandSudo) Help(ctx context.Context) (out string, err error) {
	err = c.cmd.Help()
	out = c.Out.String()
	return
}
