package local_migrations

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/system/initial/assertions"
)

type MigrationRoleNobody struct {
	adaptors *adaptors.Adaptors
}

func NewMigrationRoleNobody(adaptors *adaptors.Adaptors) *MigrationRoleNobody {
	return &MigrationRoleNobody{
		adaptors: adaptors,
	}
}

func (r *MigrationRoleNobody) Up(ctx context.Context, adaptors *adaptors.Adaptors) (err error) {
	if adaptors != nil {
		r.adaptors = adaptors
	}

	_, err = r.add(ctx)

	return
}

func (r *MigrationRoleNobody) add(ctx context.Context) (nobodyRole *m.Role, err error) {

	if nobodyRole, err = r.adaptors.Role.GetByName(ctx, "nobody"); err != nil {
		nobodyRole = &m.Role{
			Name: "nobody",
		}
		err = r.adaptors.Role.Add(ctx, nobodyRole)
		So(err, ShouldBeNil)
	}

	return
}
