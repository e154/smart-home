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
	"strings"
)

type CommandHelp struct {
	Out      *bytes.Buffer
	password *string
}

func NewCommandHelp() *CommandHelp {
	buf := new(bytes.Buffer)

	command := &CommandHelp{
		Out:      buf,
		password: new(string),
	}

	return command
}

func (c *CommandHelp) Execute(ctx context.Context, args ...string) (out string, err error) {

	var builder strings.Builder
	builder.WriteString("available commands:\n\r")
	builder.WriteString("clear\n\r")
	for cmd, _ := range commands {
		if cmd == "help" {
			continue
		}
		builder.WriteString(cmd)
		builder.WriteString("\n\r")
	}

	out = builder.String()

	return
}

func (c *CommandHelp) Help(ctx context.Context) (out string, err error) {
	return
}
