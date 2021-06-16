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
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

//go:embed swagger-ui/*
//go:embed api.swagger.json
var f embed.FS

var (
	log = common.MustGetLogger("api")
)

type Api struct {
	controllers *controllers.Controllers
	filter      *rbac.AccessFilter
}

func NewApi(controllers *controllers.Controllers,
	filter *rbac.AccessFilter) (api *Api) {
	api = &Api{controllers: controllers,
		filter: filter}
	return
}

func (a *Api) Start() error {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(a.filter.AuthInterceptor),
	)
	gw.RegisterAuthServiceServer(grpcServer, a.controllers.Auth)
	gw.RegisterStreamServiceServer(grpcServer, a.controllers.Stream)
	gw.RegisterUserServiceServer(grpcServer, a.controllers.User)
	grpc_prometheus.Register(grpcServer)

	var group errgroup.Group

	group.Go(func() error {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Error(err.Error())
		}
		return err
	})

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}))

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(50000000)),
	}

	group.Go(func() error {
		gw.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, ":3000", opts)
		gw.RegisterStreamServiceHandlerFromEndpoint(ctx, mux, ":3000", opts)
		gw.RegisterUserServiceHandlerFromEndpoint(ctx, mux, ":3000", opts)
		return nil
	})

	group.Go(func() error {
		return http.ListenAndServe(":8843", wsproxy.WebsocketProxy(mux))
	})
	group.Go(func() error {
		return http.ListenAndServe(":2662", promhttp.Handler())
	})

	swagger := http.NewServeMux()
	swagger.Handle("/", mux)
	swagger.Handle("/swagger-ui/", http.StripPrefix("/", http.FileServer(http.FS(f))))
	swagger.Handle("/api.swagger.json", http.StripPrefix("/", http.FileServer(http.FS(f))))

	go func() {
		if err = http.ListenAndServe("localhost:8080", swagger); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return group.Wait()
}
