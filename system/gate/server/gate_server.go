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

package server

import (
	"context"
	"time"

	"go.uber.org/fx"

	"github.com/e154/smart-home/common/events"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/gate/server/wsp"
)

var (
	log = logger.MustGetLogger("gate")
)

// GateServer ...
type GateServer struct {
	eventBus   bus.Bus
	cfg        Config
	proxy      *wsp.Server
	server     *Server
	gateConfig *models.GateConfig
}

// NewGateServer ...
func NewGateServer(lc fx.Lifecycle,
	eventBus bus.Bus,
	gateConfig *models.GateConfig) (gate *GateServer) {

	gate = &GateServer{
		eventBus:   eventBus,
		gateConfig: gateConfig,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return gate.Start(ctx)

		},
		OnStop: func(ctx context.Context) error {
			return gate.Shutdown(ctx)
		},
	})

	return
}

// Start ...
func (g *GateServer) Start(ctx context.Context) (err error) {

	config := &wsp.Config{
		Timeout:     time.Duration(g.gateConfig.ProxyTimeout) * time.Second,
		IdleTimeout: time.Duration(g.gateConfig.ProxyIdleTimeout) * time.Second,
		SecretKey:   g.gateConfig.ProxySecretKey,
	}
	g.proxy = wsp.NewServer(config)
	g.proxy.Start()

	cfg := &Config{
		HttpPort:  g.gateConfig.ApiHttpPort,
		HttpsPort: g.gateConfig.ApiHttpsPort,
		Debug:     g.gateConfig.ApiDebug,
		Pprof:     g.gateConfig.Pprof,
		Gzip:      g.gateConfig.ApiGzip,
		Https:     g.gateConfig.Https,
		Domain:    g.gateConfig.Domain,
	}
	g.server = NewServer(cfg, g.proxy)
	g.server.Start()

	log.Info("Started ...")

	g.eventBus.Publish("system/services/gate_server", events.EventServiceStarted{Service: "GateServer"})

	return
}

// Shutdown ...
func (g *GateServer) Shutdown(ctx context.Context) (err error) {

	log.Info("Shutdown ...")

	g.server.Shutdown(ctx)

	g.proxy.Shutdown()

	g.eventBus.Publish("system/services/gate_server", events.EventServiceStopped{Service: "GateServer"})
	return
}
