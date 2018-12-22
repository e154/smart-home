package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type Role struct {
	Name        string    `json:"name" valid:"MaxSize(254);Required"`
	Description string    `json:"description"`
	Parent      *Role     `json:"parent"`
	Children    []*Role   `json:"children"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Role) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
