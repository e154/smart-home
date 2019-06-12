package models

import (
	"github.com/e154/smart-home/common"
	"time"
	"github.com/e154/smart-home/system/validation"
)

type LogQuery struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Levels    []string   `json:"levels"`
}

type Log struct {
	Id        int64           `json:"id"`
	Body      string          `json:"body" valid:"Required"`
	Level     common.LogLevel `json:"level" valid:"Required"`
	CreatedAt time.Time       `json:"created_at"`
}

func (d *Log) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
