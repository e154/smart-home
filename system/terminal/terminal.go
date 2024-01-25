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

package terminal

import (
	"context"
	"github.com/e154/smart-home/system/jwt_manager"
	"time"

	"go.uber.org/fx"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
)

var (
	log = logger.MustGetLogger("terminal")
)

// Terminal ...
type Terminal struct {
	config     *m.AppConfig
	eventBus   bus.Bus
	adaptors   *adaptors.Adaptors
	jwtManager jwt_manager.JwtManager
	commands   map[string]ICommands
}

// NewTerminal ...
func NewTerminal(lc fx.Lifecycle,
	config *m.AppConfig,
	eventBus bus.Bus,
	adaptors *adaptors.Adaptors,
	jwtManager jwt_manager.JwtManager,
	commands map[string]ICommands) (t *Terminal) {
	t = &Terminal{
		config:     config,
		eventBus:   eventBus,
		adaptors:   adaptors,
		jwtManager: jwtManager,
		commands:   commands,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return t.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return t.Shutdown(ctx)
		},
	})

	return
}

// Start ...
func (t *Terminal) Start(_ context.Context) (err error) {
	_ = t.eventBus.Subscribe("system/terminal", t.eventHandler)
	return
}

// Shutdown ...
func (t *Terminal) Shutdown(_ context.Context) (err error) {
	_ = t.eventBus.Unsubscribe("system/terminal", t.eventHandler)
	return
}

func (t *Terminal) eventHandler(_ string, message interface{}) {
	switch v := message.(type) {
	case events.CommandTerminal:
		go t.RunCommand(v)
	}
}

func (t *Terminal) RunCommand(cmd events.CommandTerminal) {

	command := NewCommand(cmd, t.commands)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ctx = context.WithValue(ctx, "Options", CommandOptions{
		eventBus:   t.eventBus,
		adaptors:   t.adaptors,
		jwtManager: t.jwtManager,
		config:     t.config,
	})

	defer cancel()
	var messageType = "info"
	output, err := command.Execute(ctx)
	if err != nil {
		messageType = "error"
		output = err.Error()
	}
	t.eventBus.Publish("system/dashboard", events.EventDirectMessage{
		UserID:    cmd.User.Id,
		SessionID: cmd.SessionID,
		Query:     "command_response",
		Message: TerminalResponseMessage{
			Type: messageType,
			Body: output,
		},
	})
}
