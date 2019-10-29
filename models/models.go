package models

import (
	"github.com/e154/smart-home/system/validation"
	"github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("models")
)

type Validity struct {
}

func (d Validity) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
