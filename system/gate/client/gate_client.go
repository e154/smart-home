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

package client

import (
	"context"
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/system/stream"
	"github.com/google/uuid"

	"go.uber.org/fx"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/gate/client/wsp"
)

var (
	log = logger.MustGetLogger("gate")
)

// GateClient ...
type GateClient struct {
	eventBus bus.Bus
	proxy    *wsp.Client
	api      *api.Api
	stream   *stream.Stream
}

// NewGateClient ...
func NewGateClient(lc fx.Lifecycle,
	eventBus bus.Bus,
	api *api.Api,
	stream *stream.Stream) (gate *GateClient) {

	gate = &GateClient{
		eventBus: eventBus,
		api:      api,
		stream:   stream,
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

	config := &wsp.Config{
		ID:           uuid.NewString(),
		Targets:      []string{"ws://127.0.0.1:8080/gate/register"},
		PoolIdleSize: 1,
		PoolMaxSize:  100,
		SecretKey:    "",
	}

	g.proxy = wsp.NewClient(config, g.api, g.stream)
	g.proxy.Start(context.Background())

	g.eventBus.Publish("system/services/gate_client", events.EventServiceStarted{Service: "GateClient"})
	return
}

// Shutdown ...
func (g *GateClient) Shutdown() (err error) {
	g.proxy.Shutdown()
	g.eventBus.Publish("system/services/gate_client", events.EventServiceStopped{Service: "GateClient"})
	return
}
