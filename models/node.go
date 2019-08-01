package models

import (
	"github.com/e154/smart-home/system/validation"
	"time"
)

type Node struct {
	Id                int64     `json:"id"`
	Name              string    `json:"name" valid:"MaxSize(254);Required"`
	Status            string    `json:"status" valid:"Required"`
	Description       string    `json:"description"`
	Login             string    `json:"login" valid:"MaxSize(254);Required"`
	Password          string    `json:"password"`
	EncryptedPassword string    `json:"encrypted_password"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (d *Node) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}
