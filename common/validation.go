package common

import "github.com/e154/smart-home/system/validation"

type Validation struct {
}

func (d Validation) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
