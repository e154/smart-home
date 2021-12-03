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
	"github.com/e154/smart-home/api/controllers"
	gw "github.com/e154/smart-home/api/stub/api"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/rbac"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"strings"
)

//go:embed swagger-ui/*
//go:embed api.swagger.json
var f embed.FS

var (
	log = common.MustGetLogger("api")
)

// Api ...
type Api struct {
	controllers *controllers.Controllers
	filter      *rbac.AccessFilter
	cfg         Config
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
func (a *Api) Start() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lis, err := net.Listen("tcp", a.cfg.GrpcHostPort)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Infof("Serving GRPC server at %s", a.cfg.GrpcHostPort)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(a.filter.AuthInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)
	gw.RegisterAuthServiceServer(grpcServer, a.controllers.Auth)
	gw.RegisterStreamServiceServer(grpcServer, a.controllers.Stream)
	gw.RegisterUserServiceServer(grpcServer, a.controllers.User)
	gw.RegisterRoleServiceServer(grpcServer, a.controllers.Role)
	gw.RegisterScriptServiceServer(grpcServer, a.controllers.Script)
	gw.RegisterImageServiceServer(grpcServer, a.controllers.Image)
	gw.RegisterPluginServiceServer(grpcServer, a.controllers.Plugin)
	grpc_prometheus.Register(grpcServer)

	var group errgroup.Group

	group.Go(func() (err error) {
		if err = grpcServer.Serve(lis); err != nil {
			log.Error(err.Error())
		}
		return
	})

	//todo check ...
	//OrigName:     true,
	//EmitDefaults: true,
	customMarshaller := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseEnumNumbers:  false,
			EmitUnpopulated: true,
			UseProtoNames:   true,
		},
	}

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, customMarshaller),
		runtime.WithIncomingHeaderMatcher(a.CustomMatcher))

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}

	group.Go(func() error {
		gw.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		gw.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		gw.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		gw.RegisterScriptServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		gw.RegisterImageServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		gw.RegisterPluginServiceHandlerFromEndpoint(ctx, mux, a.cfg.GrpcHostPort, opts)
		return nil
	})

	if a.cfg.WsHostPort != "" {
		group.Go(func() error {
			return http.ListenAndServe(a.cfg.WsHostPort, wsproxy.WebsocketProxy(mux))
		})
		log.Infof("Serving WS server at %s", a.cfg.WsHostPort)
	}

	if a.cfg.PromHostPort != "" {
		group.Go(func() error {
			return http.ListenAndServe(a.cfg.PromHostPort, promhttp.Handler())
		})
		log.Infof("Serving PROMETHEUS server at %s", a.cfg.PromHostPort)
	}

	if a.cfg.Swagger {
		httpv1 := http.NewServeMux()
		httpv1.Handle("/", grpcHandlerFunc(grpcServer, mux))

		// upload handler
		httpv1.HandleFunc("/v1/image/upload", a.controllers.Image.MuxUploadImage())

		// uploaded and other static files
		httpv1.Handle("/upload", http.FileServer(http.Dir(common.StoragePath())))
		httpv1.Handle("/api_static", http.FileServer(http.Dir(common.StoragePath())))

		// swagger
		httpv1.Handle("/swagger-ui/", http.StripPrefix("/", http.FileServer(http.FS(f))))
		httpv1.Handle("/api.swagger.json", http.StripPrefix("/", http.FileServer(http.FS(f))))

		go func() {
			if err = http.ListenAndServe(a.cfg.HttpHostPort, httpv1); err != nil {
				log.Fatal(err.Error())
			}
		}()
		log.Infof("Serving HTTP server at %s", a.cfg.HttpHostPort)
	}

	return group.Wait()
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
