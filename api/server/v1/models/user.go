package models

import "github.com/e154/smart-home/system/validation"

type NewUserMeta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NewUserImage struct {
	Id int64 `json:"id"`
}

type NewUserRole struct {
	Name string `json:"name"`
}

type NewUser struct {
	Nickname       string         `json:"nickname" valid:"Required;MinSize(3);MaxSize(255)"`
	FirstName      string         `json:"first_name"`
	LastName       string         `json:"last_name"`
	Password       string         `json:"password" valid:"Required;MinSize(6);MaxSize(255)"`
	PasswordRepeat string         `json:"password_repeat" valid:"Required;MinSize(6);MaxSize(255)"`
	Email          string         `json:"email" valid:"Required;Email"`
	Status         string         `json:"status"`
	Lang           string         `json:"lang"`
	Avatar         *NewUserImage  `json:"avatar"`
	Role           *NewUserRole   `json:"role"`
	Meta           []*NewUserMeta `json:"meta"`
}

func (d *NewUser) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(d); !ok {
		errs = valid.Errors
	}

	return
}

type User struct {

}