package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type MapOptions struct {
	Zoom              float64 `json:"zoom"`
	ElementStateText  bool    `json:"element_state_text"`
	ElementOptionText bool    `json:"element_option_text"`
}

type Map struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name" valid:"MaxSize(254);Required"`
	Description string      `json:"description" valid:"Required"`
	Options     MapOptions  `json:"options"`
	Layers      []*MapLayer `json:"layers"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (m *Map) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
