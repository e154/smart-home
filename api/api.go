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

package api

import (
	"context"
	"embed"
	"fmt"
	"github.com/e154/smart-home/common/events"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	echopprof "github.com/hiko1129/echo-pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/e154/smart-home/api/controllers"
	gw "github.com/e154/smart-home/api/stub/api"
	publicAssets "github.com/e154/smart-home/build"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/rbac"
)

//go:embed swagger-ui/*
//go:embed api.swagger.json
var assets embed.FS

var (
	log = logger.MustGetLogger("api")
)

// Api ...
type Api struct {
	controllers *controllers.Controllers
	filter      *rbac.AccessFilter
	cfg         Config
	lis         net.Listener
	echo        *echo.Echo
	httpServer  *http.Server
	grpcServer  *grpc.Server
	eventBus    bus.Bus
}

// NewApi ...
func NewApi(controllers *controllers.Controllers,
	filter *rbac.AccessFilter,
	cfg Config,
	eventBus bus.Bus) (api *Api) {
	api = &Api{
		controllers: controllers,
		filter:      filter,
		cfg:         cfg,
		eventBus:    eventBus,
	}
	return
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// Start ...
func (a *Api) Start() (err error) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	a.lis, err = net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GrpcPort))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(a.filter.AuthInterceptor),
	)
	gw.RegisterAuthServiceServer(a.grpcServer, a.controllers.Auth)
	gw.RegisterStreamServiceServer(a.grpcServer, a.controllers.Stream)
	gw.RegisterUserServiceServer(a.grpcServer, a.controllers.User)
	gw.RegisterRoleServiceServer(a.grpcServer, a.controllers.Role)
	gw.RegisterScriptServiceServer(a.grpcServer, a.controllers.Script)
	gw.RegisterImageServiceServer(a.grpcServer, a.controllers.Image)
	gw.RegisterPluginServiceServer(a.grpcServer, a.controllers.Plugin)
	gw.RegisterZigbee2MqttServiceServer(a.grpcServer, a.controllers.Zigbee2mqtt)
	gw.RegisterEntityServiceServer(a.grpcServer, a.controllers.Entity)
	gw.RegisterAutomationServiceServer(a.grpcServer, a.controllers.Automation)
	gw.RegisterAreaServiceServer(a.grpcServer, a.controllers.Area)
	gw.RegisterDeveloperToolsServiceServer(a.grpcServer, a.controllers.Dev)
	gw.RegisterInteractServiceServer(a.grpcServer, a.controllers.Interact)
	gw.RegisterLogServiceServer(a.grpcServer, a.controllers.Logs)
	gw.RegisterDashboardServiceServer(a.grpcServer, a.controllers.Dashboard)
	gw.RegisterDashboardTabServiceServer(a.grpcServer, a.controllers.DashboardTab)
	gw.RegisterDashboardCardServiceServer(a.grpcServer, a.controllers.DashboardCard)
	gw.RegisterDashboardCardItemServiceServer(a.grpcServer, a.controllers.DashboardCardItem)
	gw.RegisterVariableServiceServer(a.grpcServer, a.controllers.Variable)
	gw.RegisterEntityStorageServiceServer(a.grpcServer, a.controllers.EntityStorage)
	gw.RegisterMetricServiceServer(a.grpcServer, a.controllers.Metric)
	gw.RegisterBackupServiceServer(a.grpcServer, a.controllers.Backup)
	gw.RegisterMessageDeliveryServiceServer(a.grpcServer, a.controllers.MessageDelivery)
	gw.RegisterActionServiceServer(a.grpcServer, a.controllers.Action)
	gw.RegisterConditionServiceServer(a.grpcServer, a.controllers.Condition)
	gw.RegisterTriggerServiceServer(a.grpcServer, a.controllers.Trigger)
	gw.RegisterMqttServiceServer(a.grpcServer, a.controllers.Mqtt)

	var group errgroup.Group

	// GRPC
	if a.cfg.GrpcPort != 0 {
		group.Go(func() (err error) {
			if err = a.grpcServer.Serve(a.lis); err != nil {
				log.Error(err.Error())
			}
			return
		})
		log.Infof("Serving GRPC server at %d", a.cfg.GrpcPort)
	}

	//todo check ...
	//OrigName:     true,
	//EmitDefaults: true,
	customMarshaller := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
			UseProtoNames:   false,
		},
	}

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, customMarshaller),
		runtime.WithIncomingHeaderMatcher(a.CustomMatcher))

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}

	grpcEndpoint := fmt.Sprintf(":%d", a.cfg.GrpcPort)
	group.Go(func() error {
		_ = gw.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterScriptServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterImageServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterPluginServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterZigbee2MqttServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterEntityServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterAutomationServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterAreaServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterDeveloperToolsServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterInteractServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterLogServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterDashboardServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterDashboardCardItemServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterDashboardCardServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterDashboardTabServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterEntityStorageServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterVariableServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterMetricServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterBackupServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterMessageDeliveryServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterActionServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterConditionServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterTriggerServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		_ = gw.RegisterMqttServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
		return nil
	})

	// HTTP
	a.echo = echo.New()

	if a.cfg.Debug {
		var format = `INFO	api/v1	[${method}] ${uri} ${status} ${latency_human} ${error}` + "\n"

		log.Info("debug enabled")
		DefaultLoggerConfig := middleware.LoggerConfig{
			Skipper:          middleware.DefaultSkipper,
			Format:           format,
			CustomTimeFormat: "2006-01-02 15:04:05.00000",
		}
		a.echo.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
		a.echo.Debug = true

		// automatically add routers for net/http/pprof
		// e.g. /debug/pprof, /debug/pprof/heap, etc.
		if a.cfg.Debug {
			log.Info("pprof enabled")
			echopprof.Wrap(a.echo)
		}
	} else {
		log.Info("debug disabled")
	}

	a.echo.HideBanner = true
	a.echo.HidePort = true

	a.registerHandlers(mux)

	go func() {
		if err := a.echo.Start(a.cfg.String()); err != nil {
			if err.Error() != "http: Server closed" {
				log.Error(err.Error())
			}
		}
	}()
	log.Infof("server started at %s", a.cfg.String())

	a.eventBus.Publish("system/services/api", events.EventServiceStarted{Service: "Api"})

	return group.Wait()
}

// Shutdown ...
func (a *Api) Shutdown(ctx context.Context) (err error) {
	if a.echo != nil {
		err = a.echo.Shutdown(ctx)
	}
	if a.grpcServer != nil {
		a.grpcServer.Stop()
	}
	if a.lis != nil {
		err = a.lis.Close()
	}

	a.eventBus.Publish("system/services/api", events.EventServiceStopped{Service: "Api"})

	return
}

// CustomMatcher ...
func (a *Api) CustomMatcher(key string) (string, bool) {
	switch key {
	case "X-Api-Key":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

func (a *Api) registerHandlers(mux *runtime.ServeMux) {

	// Swagger
	if a.cfg.Swagger {
		var contentHandler = echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(assets))))
		a.echo.GET("/swagger-ui/*", contentHandler)
		a.echo.GET("/api.swagger.json", contentHandler)
	}

	// Base
	a.echo.GET("/", echo.WrapHandler(a.controllers.Index.Index(publicAssets.F)))
	a.echo.Any("/ws", echo.WrapHandler(wsproxy.WebsocketProxy(mux)))
	a.echo.Any("/v1/*", echo.WrapHandler(grpcHandlerFunc(a.grpcServer, mux)))
	a.echo.Any("/v1/image/upload", echo.WrapHandler(a.controllers.Image.MuxUploadImage()))

	// uploaded and other static files
	fileServer := http.FileServer(http.Dir("./data/file_storage"))
	a.echo.Any("/upload/", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/upload/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		fileServer.ServeHTTP(w, r)
	})))
	staticServer := http.FileServer(http.Dir("./data/static"))
	a.echo.Any("/static/", echo.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/static/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		staticServer.ServeHTTP(w, r)
	})))

	// public
	a.echo.GET("/public/", echo.WrapHandler(http.StripPrefix("/", http.FileServer(http.FS(publicAssets.F)))))

	// media
	a.echo.Any("/stream/:entity_id/mse", a.controllers.Media.StreamMSE)

	// Cors
	a.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodHead},
	}))

}
