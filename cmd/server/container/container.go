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

package container

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/server"
	controllersV1 "github.com/e154/smart-home/api/server/v1/controllers"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/automation"
	"github.com/e154/smart-home/system/backup"
	"github.com/e154/smart-home/system/config"
	"github.com/e154/smart-home/system/entity_manager"
	"github.com/e154/smart-home/system/event_bus"
	"github.com/e154/smart-home/system/gate_client"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/e154/smart-home/system/metrics"
	"github.com/e154/smart-home/system/migrations"
	"github.com/e154/smart-home/system/mqtt"
	"github.com/e154/smart-home/system/mqtt_authenticator"
	"github.com/e154/smart-home/system/notify"
	"github.com/e154/smart-home/system/orm"
	"github.com/e154/smart-home/system/plugins"
	"github.com/e154/smart-home/system/rbac"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/storage"
	"github.com/e154/smart-home/system/stream"
	"github.com/e154/smart-home/system/zigbee2mqtt"
	"go.uber.org/fx"
)

// BuildContainer...
func BuildContainer(opt fx.Option) (app *fx.App) {

	app = fx.New(
		fx.Provide(
			config.ReadConfig,
			NewOrmConfig,
			orm.NewOrm,
			NewMigrationsConfig,
			migrations.NewMigrations,
			adaptors.NewAdaptors,
			logging.NewLogger,
			logging.NewLogDbSaver,
			scripts.NewScriptService,
			initial.NewInitial,
			NewMqttConfig,
			mqtt_authenticator.NewAuthenticator,
			mqtt.NewMqtt,
			access_list.NewAccessListService,
			rbac.NewAccessFilter,
			NewMetricConfig,
			metrics.NewMetricManager,
			NewZigbee2mqttConfig,
			zigbee2mqtt.NewZigbee2mqtt,
			storage.NewStorage,
			plugins.NewPluginManager,
			entity_manager.NewEntityManager,
			automation.NewAutomation,
			event_bus.NewEventBus,
			endpoint.NewEndpoint,
			notify.NewNotify,
			NewApiServerConfig,
			controllersV1.NewControllersV1,
			server.NewServer,
			stream.NewStreamService,
			stream.NewHub,
			NewBackupConfig,
			backup.NewBackup,
			gate_client.NewGateClient,
		),
		fx.Logger(NewPrinter()),
		opt,
	)

	return
}
