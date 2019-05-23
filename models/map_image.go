package models

import (
	"github.com/e154/smart-home/system/validation"
)

type MapImage struct {
	Id      int64  `json:"id"`
	Image   *Image `json:"image"`
	ImageId int64  `json:"image_id" valid:"Required"`
	Style   string `json:"style"`
}

func (m *MapImage) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
