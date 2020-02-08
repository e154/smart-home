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
	"github.com/e154/smart-home/common"
)

func (s *MobileServer) setControllers() {

	r := s.engine

	r.Static("/upload", common.StoragePath())

	basePath := r.Group("/api")

	v1 := basePath.Group("/v1")

	// ws
	v1.GET("/ws", s.af.Auth, s.streamService.Ws)
	v1.GET("/ws/*any", s.af.Auth, s.streamService.Ws)

	// auth
	v1.POST("/signin", s.ControllersV1.Auth.SignIn)
	v1.POST("/signout", s.af.Auth, s.ControllersV1.Auth.SignOut)
	v1.POST("/recovery", s.ControllersV1.Auth.Recovery)
	v1.POST("/reset", s.ControllersV1.Auth.Reset)
	v1.GET("/access_list", s.af.Auth, s.ControllersV1.Auth.AccessList)

	// map
	v1.GET("/map/active_elements", s.af.Auth, s.ControllersV1.Map.GetActiveElements)

	// mobile
	v1.GET("/workflows", s.af.Auth, s.ControllersV1.Workflow.GetList)
	v1.GET("/workflow/:id", s.af.Auth, s.ControllersV1.Workflow.GetById)
	v1.PUT("/workflow/:id/update_scenario", s.af.Auth, s.ControllersV1.Workflow.UpdateScenario)
}
