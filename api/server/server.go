// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/rbac"
	"github.com/e154/smart-home/system/stream"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	graceful      *graceful_service.GracefulService
	logger        *ServerLogger
	af            *rbac.AccessFilter
	streamService *stream.StreamService
	core          *core.Core
	isStarted     bool
}

// Start ...
func (s *Server) Start() {

	if s.isStarted {
		return
	}

	s.isStarted = true

	s.server = &http.Server{
		Addr:    s.Config.String(),
		Handler: s.engine,
	}

	go func() {
		// service connections
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err.Error())
		}
	}()

	log.Infof("Serving server at %s", s.Config.String())

	go s.core.Run()
}

// Shutdown ...
func (s *Server) Shutdown() {

	if !s.isStarted {
		return
	}
	s.isStarted = false

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
	log.Info("Server exiting")
}

// GetEngine ...
func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

// NewServer ...
func NewServer(cfg *Config,
	ctrls *controllers.ControllersV1,
	graceful *graceful_service.GracefulService,
	accessFilter *rbac.AccessFilter,
	streamService *stream.StreamService,
	core *core.Core) (newServer *Server) {

	logger := NewLogger()

	gin.DisableConsoleColor()
	gin.DefaultWriter = logger
	gin.DefaultErrorWriter = logger
	//if cfg.RunMode == config.ReleaseMode {
	gin.SetMode(gin.ReleaseMode)
	//} else {
	//	gin.SetMode(gin.DebugMode)
	//}

	engine := gin.New()
	engine.Use(gin.Recovery())
	//cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	engine.Use(cors.New(corsConfig))

	newServer = &Server{
		Config:        cfg,
		ControllersV1: ctrls,
		engine:        engine,
		graceful:      graceful,
		logger:        logger,
		af:            accessFilter,
		streamService: streamService,
		core:          core,
	}

	newServer.graceful.Subscribe(newServer)

	newServer.setControllers()

	return
}
