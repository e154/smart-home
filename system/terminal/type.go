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
	"context"
	"github.com/e154/bus"
	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/spf13/cobra"
)

type TerminalResponseMessage struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type UpdateAccessTokenMessage struct {
	AccessToken string `json:"access_token"`
}

var commands = map[string]ICommands{
	"sudo": NewCommandSudo(),
	"help": NewCommandHelp(),
}

func GetTerminalCommands() map[string]ICommands {
	return commands
}

type ICommands interface {
	Execute(ctx context.Context, args ...string) (out string, err error)
	Help(ctx context.Context) (out string, err error)
}

func emptyRun(*cobra.Command, []string) {}

type CommandOptions struct {
	eventBus   bus.Bus
	adaptors   *adaptors.Adaptors
	jwtManager jwt_manager.JwtManager
	config     *m.AppConfig
}
