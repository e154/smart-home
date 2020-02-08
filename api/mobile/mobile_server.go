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

package mobile

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/api/mobile/v1/controllers"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/rbac"
	"github.com/e154/smart-home/system/stream"
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"
	"net/http"
	"time"
)

var (
	log = logging.MustGetLogger("server")
)

type MobileServer struct {
	Config        *MobileServerConfig
	ControllersV1 *controllers.MobileControllersV1
	engine        *gin.Engine
	server        *http.Server
	logger        *MobileServerLogger
	af            *rbac.AccessFilter
	streamService *stream.StreamService
	gateClient    *gate_client.GateClient
}

func (s *MobileServer) Start() {

	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port),
		Handler: s.engine,
	}

	go func() {
		// service connections
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		s.gateClient.SetEngine(s.engine)
	}()

	log.Infof("Serving server at http://[::]:%d", s.Config.Port)
}

func (s *MobileServer) Shutdown() {
	log.Info("Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
	log.Info("MobileServer exiting")
}

func (s *MobileServer) GetEngine() *gin.Engine {
	return s.engine
}

func NewMobileServer(cfg *MobileServerConfig,
	ctrls *controllers.MobileControllersV1,
	graceful *graceful_service.GracefulService,
	accessFilter *rbac.AccessFilter,
	streamService *stream.StreamService,
	gateClient *gate_client.GateClient) (newServer *MobileServer) {

	logger := &MobileServerLogger{log}

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

	newServer = &MobileServer{
		Config:        cfg,
		ControllersV1: ctrls,
		engine:        engine,
		logger:        logger,
		af:            accessFilter,
		streamService: streamService,
		gateClient:    gateClient,
	}

	graceful.Subscribe(newServer)

	newServer.setControllers()

	return
}
