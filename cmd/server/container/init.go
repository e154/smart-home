package container

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/endpoint"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/initial/local_migrations"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/system/validation"
)

func MigrationList(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	validation *validation.Validate,
	scriptService scripts.ScriptService,
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
	}
}
