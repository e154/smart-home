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

package container

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/controllers"
	"github.com/e154/smart-home/common/web"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/bus"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/logging_db"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/orm"
	plugins2 "github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/system/scheduler"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/storage"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/stream/handlers"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"go.uber.org/dig"
	"go.uber.org/fx"
)

// BuildContainer ...
func BuildContainer() (container *dig.Container) {

	container = dig.New()
	_ = container.Provide(NewOrmConfig)
	_ = container.Provide(web.New)
	_ = container.Provide(validation.NewValidate)
	_ = container.Provide(orm.NewOrm)
	_ = container.Provide(NewMigrationsConfig)
	_ = container.Provide(migrations.NewMigrations)
	_ = container.Provide(adaptors.NewAdaptors)
	_ = container.Provide(scheduler.NewScheduler)
	_ = container.Provide(scripts.NewScriptService)
	_ = container.Provide(initial.NewInitial)
	_ = container.Provide(NewBackupConfig)
	_ = container.Provide(backup.NewBackup)
	_ = container.Provide(NewMqtt)
	_ = container.Provide(NewMqttCli)
	_ = container.Provide(mqtt_authenticator.NewAuthenticator)
	_ = container.Provide(access_list.NewAccessListService)
	_ = container.Provide(stream.NewStreamService)
	_ = container.Provide(gate_client.NewGateClient)
	_ = container.Provide(NewZigbee2mqttConfig)
	_ = container.Provide(zigbee2mqtt.NewZigbee2mqtt)
	_ = container.Provide(logging.NewLogger)
	_ = container.Provide(logging_db.NewLogDbSaver)
	_ = container.Provide(storage.NewStorage)
	_ = container.Provide(plugins2.NewPluginManager)
	_ = container.Provide(entity_manager.NewEntityManager)
	_ = container.Provide(automation.NewAutomation)
	_ = container.Provide(bus.NewBus)
	_ = container.Provide(endpoint.NewCommonEndpoint)
	_ = container.Provide(endpoint.NewEndpoint)
	_ = container.Provide(handlers.NewEventHandler)
	_ = container.Provide(controllers.NewControllers)
	_ = container.Provide(NewDialer)
	_ = container.Provide(jwt_manager.NewJwtManager)

	_ = container.Provide(func() (conf *models.AppConfig, err error) {
		conf, err = config.ReadConfig("conf", "config.json", "")
		conf.PgName = "smart_home_test"
		conf.Logging = false
		return
	})

	_ = container.Provide(func() (lc fx.Lifecycle) {
		return &FxNull{}
	})

	return
}
