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
}
