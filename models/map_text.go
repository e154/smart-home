package models

import (
	"github.com/e154/smart-home/system/validation"
)

type MapText struct {
	Id        int64     `json:"id"`
	Text      string    `json:"text" valid:"Required"`
	Style     string    `json:"style"`
}

func (m *MapText) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
