package models

import (
	"time"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
)

type Script struct {
	Id          int64      `gorm:"primary_key"`
	Lang        ScriptLang `json:"lang" valid:"Required"`
	Name        string     `json:"name" valid:"MaxSize(254);Required"`
	Source      string     `json:"source"`
	Description string     `json:"description"`
	Compiled    string     `json:"compiled"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (d *Script) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}