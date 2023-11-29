// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package gate_server

import (
	"context"
	"go.uber.org/fx"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/bus"
)

var (
	log = logger.MustGetLogger("gate")
)

// GateServer ...
type GateServer struct {
	eventBus bus.Bus
	server *Server
}

// NewGateServer ...
func NewGateServer(lc fx.Lifecycle,
	eventBus bus.Bus) (gate *GateServer) {

	gate = &GateServer{
		eventBus: eventBus,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return gate.Start()

		},
		OnStop: func(ctx context.Context) error {
			return gate.Shutdown()
		},
	})

	return
}

// Start ...
func (g *GateServer) Start() (err error) {

	config := &Config{
		Host:        "127.0.0.1",
		Port:        8080,
		Timeout:     1000,
		IdleTimeout: 60000,
		SecretKey:   "",
	}
	g.server = NewServer(config)
	g.server.Start()

	log.Info("Started ...")

	g.eventBus.Publish("system/services/gate_server", events.EventServiceStarted{Service: "GateServer"})

	return
}

// Shutdown ...
func (g *GateServer) Shutdown() (err error) {

	log.Info("Shutdown ...")

	g.server.Shutdown()

	g.eventBus.Publish("system/services/gate_server", events.EventServiceStopped{Service: "GateServer"})
	return
}
