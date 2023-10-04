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

package gate_client

import (
	"context"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/system/bus"
	"go.uber.org/fx"
)

// GateClient ...
type GateClient struct {
	eventBus bus.Bus
}

// NewGateClient ...
func NewGateClient(lc fx.Lifecycle,
	eventBus bus.Bus) (gate *GateClient) {

	gate = &GateClient{
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
func (g *GateClient) Start() (err error) {
	g.eventBus.Publish("system/services/gate_client", events.EventServiceStarted{Service: "GateClient"})
	return
}

// Shutdown ...
func (g *GateClient) Shutdown() (err error) {
	g.eventBus.Publish("system/services/gate_client", events.EventServiceStopped{Service: "GateClient"})
	return
}
