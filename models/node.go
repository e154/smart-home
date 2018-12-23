package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type Node struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name" valid:"MaxSize(254);Required"`
	Ip          string    `json:"ip" valid:"IP;Required"` // Must be a valid IPv4 address
	Port        int       `json:"port" valid:"Range(1, 65535);Required"`
	Status      string    `json:"status" valid:"Required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (d *Node) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
