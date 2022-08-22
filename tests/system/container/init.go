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
	return []local_migrations.Migration{}
}
