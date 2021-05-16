// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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
	"github.com/e154/smart-home/api/server/v1/controllers"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/rbac"
	"github.com/e154/smart-home/system/stream"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"time"
)

var (
	log = common.MustGetLogger("api.server")
)

// Server ...
type Server struct {
	Config        *Config
	ControllersV1 *controllers.ControllersV1
	engine        *gin.Engine
	server        *http.Server
	logger        *ServerLogger
	af            *rbac.AccessFilter
	streamService *stream.StreamService
}

// Start ...
func (s *Server) Start() (err error) {

	s.server = &http.Server{
		Addr:    s.Config.String(),
		Handler: s.engine,
	}

	go func() {
		// service connections
		if err = s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()

	log.Infof("Serving server at %s", s.Config.String())

	return
}

// Shutdown ...
func (s *Server) Shutdown() (err error) {

	if s.server == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err = s.server.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
	log.Info("Server exiting")

	return
}

// GetEngine ...
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// NewServer ...
func NewServer(lc fx.Lifecycle,
	cfg *Config,
	ctrls *controllers.ControllersV1,
	accessFilter *rbac.AccessFilter,
	streamService *stream.StreamService,
	gate *gate_client.GateClient) (newServer *Server) {

	logger := NewLogger()

	gin.DisableConsoleColor()
	gin.DefaultWriter = logger
	gin.DefaultErrorWriter = logger
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))

	// gate
	gate.SetMobileApiEngine(engine)

	newServer = &Server{
		Config:        cfg,
		ControllersV1: ctrls,
		engine:        engine,
		logger:        logger,
		af:            accessFilter,
		streamService: streamService,
	}

	newServer.setControllers()

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return newServer.Shutdown()
		},
	})

	return
}
