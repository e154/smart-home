// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/initial/local_migrations"
	"github.com/e154/smart-home/system/orm"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/validation"
)

func MigrationList(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	validation *validation.Validate,
	scriptService scripts.ScriptService,
	orm *orm.Orm,
	endpoint *endpoint.Endpoint) []local_migrations.Migration {
	return []local_migrations.Migration{
		local_migrations.NewMigrationImages(adaptors, "./"),
		local_migrations.NewMigrationTemplates(adaptors),
		local_migrations.NewMigrationAreas(adaptors),
		local_migrations.NewMigrationPlugins(adaptors),
		local_migrations.NewMigrationRoles(adaptors, accessList, validation),
		local_migrations.NewMigrationRoleNobody(adaptors),
		local_migrations.NewMigrationScripts(adaptors, scriptService),
		local_migrations.NewMigrationEntity(adaptors, endpoint),
		local_migrations.NewMigrationAutomations(adaptors, endpoint),
		local_migrations.NewMigrationDashboard(adaptors),
		local_migrations.NewMigrationEncryptor(adaptors),
		local_migrations.NewMigrationJavascript(adaptors),
		local_migrations.NewMigrationMqttBridge(adaptors),
		local_migrations.NewMigrationScheduler(adaptors),
		local_migrations.NewMigrationJavascriptV2(adaptors),
		local_migrations.NewMigrationTimezone(adaptors),
		local_migrations.NewMigrationSpeedtest(adaptors),
		local_migrations.NewMigrationBackup(adaptors),
		local_migrations.NewMigrationGate(adaptors),
		local_migrations.NewMigrationAddVar1(adaptors),
		local_migrations.NewMigrationUpdatePermissions(adaptors, accessList, orm),
		local_migrations.NewMigrationWebdav(adaptors),
	}
}
