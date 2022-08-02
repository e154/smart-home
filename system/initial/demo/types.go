package demo

import (
	"context"

	"github.com/e154/smart-home/adaptors"
)

// Demo ...
type Demo interface {
	Install(context.Context, *adaptors.Adaptors) error
}
