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

package api

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/mobile"
	mobileControllers "github.com/e154/smart-home/api/mobile/v1/controllers"
	"github.com/e154/smart-home/api/server"
	"github.com/e154/smart-home/api/server/v1/controllers"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/alexa"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/graceful_service"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/orm"
	"github.com/e154/smart-home/system/rbac"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"go.uber.org/dig"
)

// BuildContainer ...
func BuildContainer() (container *dig.Container) {

	container = dig.New()
	container.Provide(server.NewServer)
	container.Provide(server.NewServerConfig)
	container.Provide(mobile.NewMobileServer)
	container.Provide(mobile.NewMobileServerConfig)
	container.Provide(mobileControllers.NewMobileControllersV1)
	container.Provide(controllers.NewControllersV1)
	//container.Provide(config.ReadConfig)
	container.Provide(graceful_service.NewGracefulService)
	container.Provide(graceful_service.NewGracefulServicePool)
	container.Provide(graceful_service.NewGracefulServiceConfig)
	container.Provide(orm.NewOrm)
	container.Provide(orm.NewOrmConfig)
	container.Provide(core.NewCore)
	container.Provide(migrations.NewMigrations)
	container.Provide(migrations.NewMigrationsConfig)
	container.Provide(adaptors.NewAdaptors)
	container.Provide(scripts.NewScriptService)
	container.Provide(core.NewCron)
	container.Provide(initial.NewInitialService)
	container.Provide(backup.NewBackupConfig)
	container.Provide(backup.NewBackup)
	container.Provide(mqtt.NewMqtt)
	container.Provide(mqtt.NewMqttConfig)
	container.Provide(mqtt_authenticator.NewAuthenticator)
	container.Provide(access_list.NewAccessListService)
	container.Provide(rbac.NewAccessFilter)
	container.Provide(stream.NewStreamService)
	container.Provide(stream.NewHub)
	container.Provide(endpoint.NewEndpoint)
	container.Provide(gate_client.NewGateClient)
	container.Provide(notify.NewNotify)
	container.Provide(metrics.NewMetricManager)
	container.Provide(metrics.NewMetricConfig)
	container.Provide(zigbee2mqtt.NewZigbee2mqttConfig)
	container.Provide(zigbee2mqtt.NewZigbee2mqtt)
	container.Provide(logging.NewLogger)
	container.Provide(logging.NewLogDbSaver)
	container.Provide(alexa.NewAlexa)

	container.Provide(func() (conf *config.AppConfig, err error) {
		conf, err = config.ReadConfig()
		conf.PgName = "smart_home_test"
		conf.Logging = true
		return
	})

	return
}
