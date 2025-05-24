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
	"errors"
	"strings"

	"github.com/e154/smart-home/pkg/events"
)

type Command struct {
	cmd      events.CommandTerminal
	commands map[string]ICommands
}

func NewCommand(cmd events.CommandTerminal,
	commands map[string]ICommands) *Command {
	return &Command{
		cmd:      cmd,
		commands: commands,
	}
}

func (c *Command) Execute(ctx context.Context) (output string, err error) {
	if c.cmd.User == nil {
		return
	}

	c.cmd.Text = standardizeSpaces(c.cmd.Text)

	args := strings.Split(c.cmd.Text, " ")
	if len(args) == 0 {
		return
	}
	command, ok := c.commands[args[0]]
	if !ok {
		return
	}

	ctx = context.WithValue(ctx, "Command", c.cmd)

	if len(args) == 1 {
		args = append(args, "-h")
	}

	outCh := make(chan string)
	errCh := make(chan error)

	go func() {
		defer close(errCh)
		defer close(outCh)

		out, err := command.Execute(ctx, args[1:]...)
		if err != nil {
			errCh <- err
			return
		}
		outCh <- out
	}()

	select {
	case <-ctx.Done():
		err = errors.New("operation aborted, execution timed out")
		return
	case err = <-errCh:
	case output = <-outCh:
	}

	return
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
