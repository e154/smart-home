package server

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/api/server/v1/controllers"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/core"
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

type Server struct {
	Config        *ServerConfig
	ControllersV1 *controllers.ControllersV1
	engine        *gin.Engine
	server        *http.Server
	graceful      *graceful_service.GracefulService
	logger        *ServerLogger
	af            *rbac.AccessFilter
	streamService *stream.StreamService
	core          *core.Core
	gateClient    *gate_client.GateClient
}

func (s *Server) Start() {

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
		s.gateClient.Connect()
	}()

	log.Infof("Serving server at http://[::]:%d", s.Config.Port)

	go s.core.Run()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Error(err.Error())
	}
	log.Info("Server exiting")
}

func (s *Server) GetEngine() *gin.Engine {
	return s.engine
}

func NewServer(cfg *ServerConfig,
	ctrls *controllers.ControllersV1,
	graceful *graceful_service.GracefulService,
	accessFilter *rbac.AccessFilter,
	streamService *stream.StreamService,
	core *core.Core,
	gateClient *gate_client.GateClient) (newServer *Server) {

	logger := &ServerLogger{log}

	gin.DisableConsoleColor()
	gin.DefaultWriter = logger
	gin.DefaultErrorWriter = logger
	if cfg.RunMode == config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	newServer = &Server{
		Config:        cfg,
		ControllersV1: ctrls,
		engine:        engine,
		graceful:      graceful,
		logger:        logger,
		af:            accessFilter,
		streamService: streamService,
		core:          core,
		gateClient:    gateClient,
	}

	newServer.graceful.Subscribe(newServer)

	newServer.setControllers()

	return
}
