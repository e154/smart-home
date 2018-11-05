package controllers

import (
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("controllers")
)

type ControllerCommon struct {

}

func NewControllerCommon() *ControllerCommon {
	return &ControllerCommon{}
}

