package local_migrations

import (
	"context"

	"github.com/e154/smart-home/adaptors"
)

// Migration ...
type Migration interface {
	Up(context.Context, *adaptors.Adaptors) error
}
