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
	"github.com/e154/bus"
	"github.com/e154/smart-home/internal/adaptors"
	"github.com/e154/smart-home/internal/api"
	"github.com/e154/smart-home/internal/api/controllers"
	"github.com/e154/smart-home/internal/common/web"
	"github.com/e154/smart-home/internal/endpoint"
	"github.com/e154/smart-home/internal/system/access_list"
	"github.com/e154/smart-home/internal/system/automation"
	"github.com/e154/smart-home/internal/system/backup"
	"github.com/e154/smart-home/internal/system/gate/client"
	"github.com/e154/smart-home/internal/system/initial"
	localMigrations "github.com/e154/smart-home/internal/system/initial/local_migrations"
	"github.com/e154/smart-home/internal/system/jwt_manager"
	"github.com/e154/smart-home/internal/system/logging"
	"github.com/e154/smart-home/internal/system/logging_db"
	"github.com/e154/smart-home/internal/system/logging_ws"
	"github.com/e154/smart-home/internal/system/media"
	"github.com/e154/smart-home/internal/system/migrations"
	"github.com/e154/smart-home/internal/system/mqtt"
	"github.com/e154/smart-home/internal/system/mqtt_authenticator"
	"github.com/e154/smart-home/internal/system/orm"
	"github.com/e154/smart-home/internal/system/rbac"
	"github.com/e154/smart-home/internal/system/scheduler"
	"github.com/e154/smart-home/internal/system/scripts"
	"github.com/e154/smart-home/internal/system/storage"
	"github.com/e154/smart-home/internal/system/stream"
	"github.com/e154/smart-home/internal/system/stream/handlers"
	"github.com/e154/smart-home/internal/system/supervisor"
	"github.com/e154/smart-home/internal/system/terminal"
	"github.com/e154/smart-home/internal/system/validation"
	"github.com/e154/smart-home/internal/system/zigbee2mqtt"
	"go.uber.org/fx"
)

// BuildContainer ...
func BuildContainer(opt fx.Option) (app *fx.App) {

	app = fx.New(
		fx.Provide(
			ReadConfig,
			validation.NewValidate,
			NewOrmConfig,
			bus.NewBus,
			orm.NewOrm,
			backup.NewBackup,
			NewMigrationsConfig,
			migrations.NewMigrations,
			web.New,
			adaptors.NewAdaptors,
			scheduler.NewScheduler,
			NewLoggerConfig,
			logging.NewLogger,
			logging_db.NewLogDbSaver,
			logging_ws.NewLogWsSaver,
			terminal.GetTerminalCommands,
			terminal.NewTerminal,
			scripts.NewScriptService,
			MigrationList,
			localMigrations.NewMigrations,
			NewDemo,
			media.NewMedia,
			initial.NewInitial,
			NewMqttConfig,
			mqtt_authenticator.NewAuthenticator,
			mqtt.NewMqtt,
			access_list.NewAccessListService,
			rbac.NewEchoAccessFilter,
			jwt_manager.NewJwtManager,
			NewZigbee2mqttConfig,
			zigbee2mqtt.NewZigbee2mqtt,
			storage.NewStorage,
			supervisor.NewSupervisor,
			automation.NewAutomation,
			endpoint.NewCommonEndpoint,
			endpoint.NewEndpoint,
			NewApiConfig,
			api.NewApi,
			controllers.NewControllers,
			stream.NewStreamService,
			handlers.NewEventHandler,
			NewBackupConfig,
			client.NewGateClient,
		),
		fx.Logger(NewPrinter()),
		opt,
	)

	return
}
