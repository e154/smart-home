package mobile

import (
	"context"
	"fmt"
	"github.com/e154/smart-home/api/mobile/v1/controllers"
	"github.com/e154/smart-home/system/config"
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
		s.gateClient.Connect()
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
	if cfg.RunMode == config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

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
