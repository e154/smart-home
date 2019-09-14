package models

import "github.com/e154/smart-home/system/validation"

type MapZone struct {
	Id   int64  `json:"id"`
	Name string `json:"name" valid:"Required"`
}

func (d *MapZone) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
