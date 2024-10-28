package controllers

import (
	"github.com/e154/smart-home/pkg/logger"
)

var (
	log = logger.MustGetLogger("controllers")
)

// ControllerCommon ...
type ControllerCommon struct {
	Mode string
}

// NewControllerCommon ...
func NewControllerCommon(mode string) *ControllerCommon {
	return &ControllerCommon{
		Mode: mode,
	}
}
