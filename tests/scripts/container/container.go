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
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/jwt_manager"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/orm"
	plugins2 "github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/storage"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"go.uber.org/dig"
)

// BuildContainer ...
func BuildContainer() (container *dig.Container) {

	container = dig.New()
	container.Provide(NewOrmConfig)
	container.Provide(orm.NewOrm)
	container.Provide(NewMigrationsConfig)
	container.Provide(migrations.NewMigrations)
	container.Provide(adaptors.NewAdaptors)
	container.Provide(scripts.NewScriptService)
	container.Provide(initial.NewInitial)
	container.Provide(NewBackupConfig)
	container.Provide(backup.NewBackup)
	container.Provide(NewMqttConfig)
	container.Provide(mqtt.NewMqtt)
	container.Provide(mqtt_authenticator.NewAuthenticator)
	container.Provide(access_list.NewAccessListService)
	container.Provide(stream.NewStreamService)
	container.Provide(stream.NewHub)
	container.Provide(gate_client.NewGateClient)
	container.Provide(NewMetricConfig)
	container.Provide(metrics.NewMetricManager)
	container.Provide(NewZigbee2mqttConfig)
	container.Provide(zigbee2mqtt.NewZigbee2mqtt)
	container.Provide(logging.NewLogger)
	container.Provide(logging.NewLogDbSaver)
	container.Provide(storage.NewStorage)
	container.Provide(plugins2.NewPluginManager)
	container.Provide(entity_manager.NewEntityManager)
	container.Provide(automation.NewAutomation)
	container.Provide(event_bus.NewEventBus)
	container.Provide(endpoint.NewEndpoint)
	container.Provide(jwt_manager.NewJwtManager)

	container.Provide(func() (conf *models.AppConfig, err error) {
		conf, err = config.ReadConfig("conf", "config.json", "")()
		conf.PgName = "smart_home_test"
		conf.Logging = false
		return
	})

	return
}
