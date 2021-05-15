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

package server

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/swaggo/gin-swagger/swaggerFiles"
	"github.com/gin-gonic/gin"
)

func (s *Server) setControllers() {

	r := s.engine

	r.Static("/upload", common.StoragePath())
	r.Static("/api_static", common.StaticPath())

	basePath := r.Group("/api")

	v1 := basePath.Group("/v1")
	v1.GET("/", s.ControllersV1.Index.Index)
	v1.GET("/swagger", func(context *gin.Context) {
		context.Redirect(302, "/api/v1/swagger/index.html")
	})
	v1.GET("/swagger/*any", s.ControllersV1.Swagger.WrapHandler(swaggerFiles.Handler))

	// ws
	v1.GET("/ws", s.af.Auth, s.streamService.Ws)
	v1.GET("/ws/*any", s.af.Auth, s.streamService.Ws)

	// auth
	v1.POST("/signin", s.ControllersV1.Auth.SignIn)
	v1.POST("/signout", s.af.Auth, s.ControllersV1.Auth.SignOut)
	v1.POST("/recovery", s.ControllersV1.Auth.Recovery)
	v1.POST("/reset", s.ControllersV1.Auth.Reset)
	v1.GET("/access_list", s.af.Auth, s.ControllersV1.Auth.AccessList)

	// scripts
	v1.POST("/script", s.af.Auth, s.ControllersV1.Script.Add)
	v1.GET("/script/:id", s.af.Auth, s.ControllersV1.Script.GetById)
	v1.PUT("/script/:id", s.af.Auth, s.ControllersV1.Script.Update)
	v1.DELETE("/script/:id", s.af.Auth, s.ControllersV1.Script.Delete)
	v1.GET("/scripts", s.af.Auth, s.ControllersV1.Script.GetList)
	v1.POST("/script/:id/exec", s.af.Auth, s.ControllersV1.Script.Exec)
	v1.POST("/script/:id/copy", s.af.Auth, s.ControllersV1.Script.Copy)
	v1.POST("/script/:id/exec_src", s.af.Auth, s.ControllersV1.Script.ExecSrc)
	v1.GET("/scripts/search", s.af.Auth, s.ControllersV1.Script.Search)

	// role
	v1.POST("/role", s.af.Auth, s.ControllersV1.Role.Add)
	v1.GET("/role/:name", s.af.Auth, s.ControllersV1.Role.GetByName)
	v1.GET("/role/:name/access_list", s.af.Auth, s.ControllersV1.Role.GetAccessList)
	v1.PUT("/role/:name/access_list", s.af.Auth, s.ControllersV1.Role.UpdateAccessList)
	v1.PUT("/role/:name", s.af.Auth, s.ControllersV1.Role.Update)
	v1.DELETE("/role/:name", s.af.Auth, s.ControllersV1.Role.Delete)
	v1.GET("/roles", s.af.Auth, s.ControllersV1.Role.GetList)
	v1.GET("/roles/search", s.af.Auth, s.ControllersV1.Role.Search)

	// user
	v1.POST("/user", s.af.Auth, s.ControllersV1.User.Add)
	v1.GET("/user/:id", s.af.Auth, s.ControllersV1.User.GetById)
	v1.PUT("/user/:id", s.af.Auth, s.ControllersV1.User.Update)
	v1.DELETE("/user/:id", s.af.Auth, s.ControllersV1.User.Delete)
	v1.PUT("/user/:id/update_status", s.af.Auth, s.ControllersV1.User.UpdateStatus)
	v1.GET("/users", s.af.Auth, s.ControllersV1.User.GetList)

	// images
	v1.POST("/image", s.af.Auth, s.ControllersV1.Image.Add)
	v1.GET("/image/:id", s.af.Auth, s.ControllersV1.Image.GetById)
	v1.GET("/images", s.af.Auth, s.ControllersV1.Image.GetList)
	v1.POST("/image/upload", s.af.Auth, s.ControllersV1.Image.Upload)
	v1.PUT("/image/:id", s.af.Auth, s.ControllersV1.Image.Update)
	v1.DELETE("/image/:id", s.af.Auth, s.ControllersV1.Image.Delete)

	// logs
	v1.POST("/log", s.af.Auth, s.ControllersV1.Log.Add)
	v1.GET("/log/:id", s.af.Auth, s.ControllersV1.Log.GetById)
	v1.DELETE("/log/:id", s.af.Auth, s.ControllersV1.Log.Delete)
	v1.GET("/logs", s.af.Auth, s.ControllersV1.Log.GetList)
	v1.GET("/logs/search", s.af.Auth, s.ControllersV1.Log.Search)

	// templates
	v1.POST("/template", s.af.Auth, s.ControllersV1.Template.Add)
	v1.GET("/template/:name", s.af.Auth, s.ControllersV1.Template.GetByName)
	v1.GET("/templates", s.af.Auth, s.ControllersV1.Template.GetList)
	v1.PUT("/template/:name", s.af.Auth, s.ControllersV1.Template.Update)
	v1.DELETE("/template/:name", s.af.Auth, s.ControllersV1.Template.Delete)
	v1.GET("/templates/search", s.af.Auth, s.ControllersV1.Template.Search)
	v1.POST("/templates/preview", s.af.Auth, s.ControllersV1.Template.Preview)

	// template items
	v1.POST("/template_item", s.af.Auth, s.ControllersV1.TemplateItem.Add)
	v1.GET("/template_item/:name", s.af.Auth, s.ControllersV1.TemplateItem.GetByName)
	v1.GET("/template_items", s.af.Auth, s.ControllersV1.TemplateItem.GetList)
	v1.GET("/template_items/tree", s.af.Auth, s.ControllersV1.TemplateItem.GetTree)
	v1.PUT("/template_items/tree", s.af.Auth, s.ControllersV1.TemplateItem.UpdateTree)
	v1.PUT("/template_item/:name", s.af.Auth, s.ControllersV1.TemplateItem.Update)
	v1.PUT("/template_items/status/:name", s.af.Auth, s.ControllersV1.TemplateItem.UpdateStatus)
	v1.DELETE("/template_item/:name", s.af.Auth, s.ControllersV1.TemplateItem.Delete)

	// notify
	v1.GET("/notifr/config", s.af.Auth, s.ControllersV1.Notifr.GetSettings)
	v1.PUT("/notifr/config", s.af.Auth, s.ControllersV1.Notifr.Update)
	v1.GET("/notifrs", s.af.Auth, s.ControllersV1.Notifr.GetList)
	v1.DELETE("/notifr/:id", s.af.Auth, s.ControllersV1.Notifr.Delete)
	v1.POST("/notifr/:id/repeat", s.af.Auth, s.ControllersV1.Notifr.Repeat)
	v1.POST("/notifr", s.af.Auth, s.ControllersV1.Notifr.Send)

	// mqtt
	v1.DELETE("/mqtt/client/:id", s.af.Auth, s.ControllersV1.Mqtt.CloseClient)
	v1.GET("/mqtt/client/:id", s.af.Auth, s.ControllersV1.Mqtt.GetClientById)
	v1.GET("/mqtt/client/:id/session", s.af.Auth, s.ControllersV1.Mqtt.GetSession)
	v1.GET("/mqtt/client/:id/subscriptions", s.af.Auth, s.ControllersV1.Mqtt.GetSubscriptions)
	v1.DELETE("/mqtt/client/:id/topic", s.af.Auth, s.ControllersV1.Mqtt.Unsubscribe)
	v1.GET("/mqtt/clients", s.af.Auth, s.ControllersV1.Mqtt.GetClients)
	v1.POST("/mqtt/publish", s.af.Auth, s.ControllersV1.Mqtt.Publish)
	v1.GET("/mqtt/sessions", s.af.Auth, s.ControllersV1.Mqtt.GetSessions)
	v1.GET("/mqtt/search_topic", s.af.Auth, s.ControllersV1.Mqtt.SearchTopic)

	// version
	v1.GET("/version", s.ControllersV1.Version.Version)

	// zigbee2mqtt
	v1.POST("/zigbee2mqtt", s.af.Auth, s.ControllersV1.Zigbee2mqtt.Add)
	v1.GET("/zigbee2mqtt/:id", s.af.Auth, s.ControllersV1.Zigbee2mqtt.GetById)
	v1.PUT("/zigbee2mqtt/:id", s.af.Auth, s.ControllersV1.Zigbee2mqtt.Update)
	v1.DELETE("/zigbee2mqtt/:id", s.af.Auth, s.ControllersV1.Zigbee2mqtt.Delete)
	v1.GET("/zigbee2mqtts", s.af.Auth, s.ControllersV1.Zigbee2mqtt.GetList)
	v1.POST("/zigbee2mqtt/:id/reset", s.af.Auth, s.ControllersV1.Zigbee2mqtt.Reset)
	v1.POST("/zigbee2mqtt/:id/device_ban", s.af.Auth, s.ControllersV1.Zigbee2mqtt.DeviceBan)
	v1.POST("/zigbee2mqtt/:id/device_whitelist", s.af.Auth, s.ControllersV1.Zigbee2mqtt.DeviceWhitelist)
	v1.GET("/zigbee2mqtt/:id/networkmap", s.af.Auth, s.ControllersV1.Zigbee2mqtt.Networkmap)
	v1.POST("/zigbee2mqtt/:id/update_networkmap", s.af.Auth, s.ControllersV1.Zigbee2mqtt.UpdateNetworkmap)
	v1.PATCH("/zigbee2mqtts/device_rename", s.af.Auth, s.ControllersV1.Zigbee2mqtt.DeviceRename)
	v1.GET("/zigbee2mqtts/search_device", s.af.Auth, s.ControllersV1.Zigbee2mqtt.Search)

	// alexa
	v1.POST("/alexa", s.af.Auth, s.ControllersV1.Alexa.Add)
	v1.GET("/alexa/:id", s.af.Auth, s.ControllersV1.Alexa.GetById)
	v1.PUT("/alexa/:id", s.af.Auth, s.ControllersV1.Alexa.Update)
	v1.DELETE("/alexa/:id", s.af.Auth, s.ControllersV1.Alexa.Delete)
	v1.GET("/alexas", s.af.Auth, s.ControllersV1.Alexa.GetList)

	// entities
	v1.POST("/entity", s.af.Auth, s.ControllersV1.Entity.Add)
	v1.GET("/entity/:id", s.af.Auth, s.ControllersV1.Entity.GetById)
	v1.PUT("/entity/:id", s.af.Auth, s.ControllersV1.Entity.Update)
	v1.DELETE("/entity/:id", s.af.Auth, s.ControllersV1.Entity.Delete)
	v1.GET("/entities", s.af.Auth, s.ControllersV1.Entity.GetList)
	v1.GET("/entities/search", s.af.Auth, s.ControllersV1.Entity.Search)

	// developer tools
	v1.GET("/developer_tools/states", s.af.Auth, s.ControllersV1.DeveloperTools.GetStateList)
	v1.PATCH("/developer_tools/state", s.af.Auth, s.ControllersV1.DeveloperTools.UpdateState)
	v1.GET("/developer_tools/events", s.af.Auth, s.ControllersV1.DeveloperTools.GetEventList)
}
