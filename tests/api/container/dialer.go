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

package container

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/e154/smart-home/api/controllers"
	gw "github.com/e154/smart-home/api/stub/api"
)

// Dialer ...
type Dialer struct {
	controllers *controllers.Controllers
}

// NewDialer ...
func NewDialer(controllers *controllers.Controllers) *Dialer {
	return &Dialer{controllers: controllers}
}

// Call ...
func (d *Dialer) Call() func(context.Context, string) (net.Conn, error) {

	grpcServer := grpc.NewServer()

	gw.RegisterAuthServiceServer(grpcServer, d.controllers.Auth)
	gw.RegisterStreamServiceServer(grpcServer, d.controllers.Stream)
	gw.RegisterUserServiceServer(grpcServer, d.controllers.User)
	gw.RegisterRoleServiceServer(grpcServer, d.controllers.Role)
	gw.RegisterScriptServiceServer(grpcServer, d.controllers.Script)
	gw.RegisterImageServiceServer(grpcServer, d.controllers.Image)
	gw.RegisterPluginServiceServer(grpcServer, d.controllers.Plugin)
	gw.RegisterZigbee2MqttServiceServer(grpcServer, d.controllers.Zigbee2mqtt)
	gw.RegisterAreaServiceServer(grpcServer, d.controllers.Area)
	gw.RegisterEntityServiceServer(grpcServer, d.controllers.Entity)
	gw.RegisterAutomationServiceServer(grpcServer, d.controllers.Automation)
	gw.RegisterVariableServiceServer(grpcServer, d.controllers.Variable)

	listener := bufconn.Listen(1024 * 1024)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
