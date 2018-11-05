package main

import (
	"github.com/e154/smart-home/api/server_v1"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/dig"
	"github.com/e154/smart-home/api/server_v1/controllers"
)

func BuildContainer() (container *dig.Container) {

	container = dig.New()

	container.Provide(server.NewServer)
	container.Provide(server.NewServerConfig)
	container.Provide(controllers.NewControllers)
	container.Provide(config.ReadConfig)

	return
}
