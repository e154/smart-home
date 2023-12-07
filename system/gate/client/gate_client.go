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
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api"
	"github.com/e154/smart-home/system/stream"
	"github.com/google/uuid"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"net/url"

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
	eventBus  bus.Bus
	proxy     *wsp.Client
	api       *api.Api
	stream    *stream.Stream
	adaptors  *adaptors.Adaptors
	inProcess *atomic.Bool
}

// NewGateClient ...
func NewGateClient(lc fx.Lifecycle,
	eventBus bus.Bus,
	api *api.Api,
	stream *stream.Stream,
	adaptors *adaptors.Adaptors) (gate *GateClient) {

	gate = &GateClient{
		eventBus:  eventBus,
		api:       api,
		stream:    stream,
		adaptors:  adaptors,
		inProcess: atomic.NewBool(false),
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
	g.initWspServer()
	_ = g.eventBus.Subscribe("system/models/variables/+", g.eventHandler, false)
	return
}

// Shutdown ...
func (g *GateClient) Shutdown() (err error) {
	if g.proxy != nil {
		g.proxy.Shutdown()
	}
	_ = g.eventBus.Unsubscribe("system/models/variables/+", g.eventHandler)
	g.eventBus.Publish("system/services/gate_client", events.EventServiceStopped{Service: "GateClient"})
	return
}

func (g *GateClient) initWspServer() {
	if !g.inProcess.CompareAndSwap(false, true) {
		return
	}
	defer g.inProcess.Store(false)

	fmt.Println("initWspServer")

	if g.proxy != nil {
		g.proxy.Shutdown()
		g.proxy = nil
	}

	list, _, err := g.adaptors.Variable.List(context.Background(), 6, 0, "", "", true, "gateClientId,gateClientSecretKey,gateClientServerHost,gateClientServerPort,gateClientPoolIdleSize,gateClientPoolMaxSize")
	if err != nil {
		log.Error(err.Error())
		return
	}

	cfg := &Config{}
	for _, variable := range list {
		switch variable.Name {
		case "gateClientId":
			cfg.Id = variable.Value
		case "gateClientSecretKey":
			cfg.SecretKey = variable.Value
		case "gateClientServerHost":
			cfg.ServerHost = variable.Value
		case "gateClientServerPort":
			cfg.ServerPort = variable.GetInt()
		case "gateClientPoolIdleSize":
			cfg.PoolIdleSize = variable.GetInt()
		case "gateClientPoolMaxSize":
			cfg.PoolMaxSize = variable.GetInt()
		}
	}

	if cfg.ServerHost == "" || cfg.ServerPort == 0 {
		return
	}

	if cfg.Id == "" {
		cfg.Id = uuid.NewString()
	}

	if cfg.PoolIdleSize == 0 {
		cfg.PoolIdleSize = 1
	}

	if cfg.PoolMaxSize == 0 {
		cfg.PoolMaxSize = 100
	}

	uri := url.URL{
		Scheme: "ws",
		Host:   fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort),
		Path:   "gate/register",
	}

	config := &wsp.Config{
		ID:           cfg.Id,
		Targets:      []string{uri.String()},
		PoolIdleSize: cfg.PoolIdleSize,
		PoolMaxSize:  cfg.PoolMaxSize,
		SecretKey:    cfg.SecretKey,
	}

	g.proxy = wsp.NewClient(config, g.api, g.stream)
	g.proxy.Start(context.Background())
}

func (g *GateClient) eventHandler(_ string, message interface{}) {

	switch v := message.(type) {
	case events.EventUpdatedVariableModel:
		switch v.Name {
		case "gateClientId", "gateClientSecretKey", "gateClientServerHost", "gateClientServerPort",
			"gateClientPoolIdleSize", "gateClientPoolMaxSize":
			log.Infof("updated settings name %s", v.Name)
			g.initWspServer()
		}
	}
}
