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
	"net"
	"net/http"
	"strings"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/e154/smart-home/api/controllers"
	gw "github.com/e154/smart-home/api/stub/api"
	publicAssets "github.com/e154/smart-home/build"
	"github.com/e154/smart-home/common/logger"
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
	httpServer  *http.Server
	promServer  *http.Server
	grpcServer  *grpc.Server
}

// NewApi ...
func NewApi(controllers *controllers.Controllers,
	filter *rbac.AccessFilter,
	cfg Config) (api *Api) {
	api = &Api{
		controllers: controllers,
		filter:      filter,
		cfg:         cfg,
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

	a.lis, err = net.Listen("tcp", a.cfg.GrpcHostPort)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	a.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(a.filter.AuthInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
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
	grpc_prometheus.Register(a.grpcServer)

	var group errgroup.Group

	// GRPC
	if a.cfg.GrpcHostPort != "" {
		group.Go(func() (err error) {
			if err = a.grpcServer.Serve(a.lis); err != nil {
				log.Error(err.Error())
			}
			return
		})
		log.Infof("Serving GRPC server at %s", a.cfg.GrpcHostPort)
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

	group.Go(func() error {
		_ = gw.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterScriptServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterImageServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterPluginServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterZigbee2MqttServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterEntityServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterAutomationServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterAreaServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterDeveloperToolsServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterInteractServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterLogServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterDashboardServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterDashboardCardItemServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterDashboardCardServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterDashboardTabServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterEntityStorageServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterVariableServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterMetricServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterBackupServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		_ = gw.RegisterMessageDeliveryServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		return nil
	})

	// PROMETHEUS
	if a.cfg.PromHostPort != "" {
		a.promServer = &http.Server{Addr: a.cfg.PromHostPort, Handler: promhttp.Handler()}
		group.Go(func() error {
			return a.promServer.ListenAndServe()
		})
		log.Infof("Serving PROMETHEUS server at %s", a.cfg.PromHostPort)
	}

	// HTTP
	httpv1 := http.NewServeMux()
	httpv1.HandleFunc("/", a.controllers.Index.Index(publicAssets.F))
	httpv1.Handle("/v1/", grpcHandlerFunc(a.grpcServer, mux))
	httpv1.Handle("/ws", wsproxy.WebsocketProxy(mux))
	httpv1.HandleFunc("/v1/image/upload", a.controllers.Image.MuxUploadImage())

	// uploaded and other static files
	fileServer := http.FileServer(http.Dir("./data/file_storage"))
	httpv1.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/upload/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		fileServer.ServeHTTP(w, r)
	})

	staticServer := http.FileServer(http.Dir("./data/static"))
	httpv1.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		r.RequestURI = strings.ReplaceAll(r.RequestURI, "/static/", "/")
		r.URL, _ = r.URL.Parse(r.RequestURI)
		staticServer.ServeHTTP(w, r)
	})

	// public
	httpv1.Handle("/public/", http.StripPrefix("/", http.FileServer(http.FS(publicAssets.F))))

	// swagger
	if a.cfg.Swagger {
		httpv1.Handle("/swagger-ui/", http.StripPrefix("/", http.FileServer(http.FS(assets))))
		httpv1.Handle("/api.swagger.json", http.StripPrefix("/", http.FileServer(http.FS(assets))))
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		Debug:            a.cfg.Debug,
	})

	a.httpServer = &http.Server{Addr: a.cfg.HttpHostPort, Handler: cors.Handler(httpv1)}
	group.Go(func() error {
		return a.httpServer.ListenAndServe()
	})
	log.Infof("Serving HTTP server at %s", a.cfg.HttpHostPort)

	return group.Wait()
}

// Shutdown ...
func (a *Api) Shutdown(ctx context.Context) (err error) {
	if a.httpServer != nil {
		err = a.httpServer.Shutdown(ctx)
	}
	if a.promServer != nil {
		err = a.promServer.Shutdown(ctx)
	}
	if a.grpcServer != nil {
		a.grpcServer.Stop()
	}
	if a.lis != nil {
		err = a.lis.Close()
	}
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
