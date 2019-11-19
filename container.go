package main

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/mobile"
	mobileControllers "github.com/e154/smart-home/api/mobile/v1/controllers"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/server/v1/controllers"
	"github.com/e154/smart-home/api/websocket"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/dig"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/orm"
	"github.com/e154/smart-home/system/rbac"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/services"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/telemetry"
)

func BuildContainer() (container *dig.Container) {

	container = dig.New()
	container.Provide(server.NewServer)
	container.Provide(server.NewServerConfig)
	container.Provide(mobile.NewMobileServer)
	container.Provide(mobile.NewMobileServerConfig)
	container.Provide(mobileControllers.NewMobileControllersV1)
	container.Provide(websocket.NewWebSocket)
	container.Provide(controllers.NewControllersV1)
	container.Provide(config.ReadConfig)
	container.Provide(graceful_service.NewGracefulService)
	container.Provide(graceful_service.NewGracefulServicePool)
	container.Provide(graceful_service.NewGracefulServiceConfig)
	container.Provide(orm.NewOrm)
	container.Provide(orm.NewOrmConfig)
	container.Provide(core.NewCore)
	container.Provide(migrations.NewMigrations)
	container.Provide(migrations.NewMigrationsConfig)
	container.Provide(adaptors.NewAdaptors)
	container.Provide(logging.NewLogrus)
	container.Provide(scripts.NewScriptService)
	container.Provide(core.NewCron)
	container.Provide(initial.NewInitialService)
	container.Provide(backup.NewBackupConfig)
	container.Provide(backup.NewBackup)
	container.Provide(services.NewServices)
	container.Provide(mqtt.NewMqtt)
	container.Provide(mqtt.NewMqttConfig)
	container.Provide(mqtt.NewAuthenticator)
	container.Provide(access_list.NewAccessListService)
	container.Provide(rbac.NewAccessFilter)
	container.Provide(stream.NewStreamService)
	container.Provide(stream.NewHub)
	container.Provide(telemetry.NewTelemetry)
	container.Provide(logging.NewLogBackend)
	container.Provide(logging.NewLogDbSaver)
	container.Provide(endpoint.NewEndpoint)
	container.Provide(gate_client.NewGateClient)
	container.Provide(notify.NewNotify)

	return
}
