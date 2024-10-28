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
	"github.com/e154/smart-home/internal/endpoint"
	"github.com/e154/smart-home/internal/system/access_list"
	local_migrations2 "github.com/e154/smart-home/internal/system/initial/local_migrations"
	"github.com/e154/smart-home/internal/system/orm"
	"github.com/e154/smart-home/internal/system/validation"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/scripts"
)

func MigrationList(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	validation *validation.Validate,
	scriptService scripts.ScriptService,
	orm *orm.Orm,
	endpoint *endpoint.Endpoint) []local_migrations2.Migration {
	return []local_migrations2.Migration{
		local_migrations2.NewMigrationInit(adaptors),
		local_migrations2.NewMigrationImages(adaptors, "./"),
		local_migrations2.NewMigrationTemplates(adaptors),
		local_migrations2.NewMigrationAreas(adaptors),
		local_migrations2.NewMigrationPlugins(adaptors),
		local_migrations2.NewMigrationRoles(adaptors, accessList, validation),
		local_migrations2.NewMigrationRoleNobody(adaptors),
		local_migrations2.NewMigrationScripts(adaptors, scriptService),
		local_migrations2.NewMigrationEntity(adaptors, endpoint),
		local_migrations2.NewMigrationAutomations(adaptors, endpoint),
		local_migrations2.NewMigrationDashboard(adaptors),
		local_migrations2.NewMigrationEncryptor(adaptors),
		local_migrations2.NewMigrationJavascript(adaptors),
		local_migrations2.NewMigrationMqttBridge(adaptors),
		local_migrations2.NewMigrationScheduler(adaptors),
		local_migrations2.NewMigrationJavascriptV2(adaptors),
		local_migrations2.NewMigrationTimezone(adaptors),
		local_migrations2.NewMigrationSpeedtest(adaptors),
		local_migrations2.NewMigrationBackup(adaptors),
		local_migrations2.NewMigrationGate(adaptors),
		local_migrations2.NewMigrationAddVar1(adaptors),
		local_migrations2.NewMigrationUpdatePermissions(adaptors, accessList, orm),
		local_migrations2.NewMigrationWebdav(adaptors),
		local_migrations2.NewMigrationAddVar2(adaptors),
		local_migrations2.NewMigrationAutocert(adaptors),
		local_migrations2.NewMigrationPachka(adaptors),
		local_migrations2.NewMigrationWebhook(adaptors),
		local_migrations2.NewMigrationRemoveTriggersPlugin(adaptors),
		local_migrations2.NewMigrationMdns(adaptors),
	}
}
