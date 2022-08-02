package container

import (
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/initial/local_migrations"
	"github.com/e154/smart-home/system/validation"
)

func MigrationList(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	validation *validation.Validate) []local_migrations.Migration {
	return []local_migrations.Migration{
		local_migrations.NewMigrationImages(adaptors, "data"),
		local_migrations.NewMigrationRoles(adaptors, accessList, validation),
		local_migrations.NewMigrationTemplates(adaptors),
		local_migrations.NewMigrationAreas(adaptors),
		local_migrations.NewMigrationPlugins(adaptors),
		local_migrations.NewMigrationZones(adaptors),
	}
}
