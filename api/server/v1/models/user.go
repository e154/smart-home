package models

import (
	"github.com/e154/smart-home/system/validation"
	"time"
)

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

type UserByIdModelParent struct {
	Id       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

type UserByIdModelMeta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UserByIdModelUserHistory struct {
	Ip   string    `json:"ip"`
	Time time.Time `json:"time"`
}

type UserByIdModel struct {
	Id                  int64                       `json:"id"`
	Nickname            string                      `json:"nickname" valid:"Required;MinSize(3);MaxSize(255)"`
	FirstName           string                      `json:"first_name" valid:"MaxSize(255)"`
	LastName            string                      `json:"last_name" valid:"MaxSize(255)"`
	Email               string                      `json:"email" valid:"Required;Email"`
	Status              string                      `json:"status" valid:"MaxSize(255)"`
	AuthenticationToken string                      `json:"authentication_token"`
	Image               *Image                      `json:"image"`
	SignInCount         int64                       `json:"sign_in_count"`
	CurrentSignInIp     string                      `json:"current_sign_in_ip"`
	LastSignInIp        string                      `json:"last_sign_in_ip"`
	Lang                string                      `json:"lang"`
	User                *UserByIdModelParent        `json:"user"`
	Role                *Role                       `json:"role"`
	RoleName            string                      `json:"role_name"`
	Meta                []*UserByIdModelMeta        `json:"meta"`
	ResetPasswordSentAt *time.Time                  `json:"reset_password_sent_at"`
	CurrentSignInAt     *time.Time                  `json:"current_sign_in_at"`
	LastSignInAt        *time.Time                  `json:"last_sign_in_at"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
	DeletedAt           *time.Time                  `json:"deleted_at"`
	History             []*UserByIdModelUserHistory `json:"history"`
}

type UserShotModel struct {
	Id        int64                `json:"id"`
	Nickname  string               `json:"nickname" valid:"Required;MinSize(3);MaxSize(255)"`
	FirstName string               `json:"first_name" valid:"MaxSize(255)"`
	LastName  string               `json:"last_name" valid:"MaxSize(255)"`
	Email     string               `json:"email" valid:"Required;Email"`
	Status    string               `json:"status" valid:"MaxSize(255)"`
	Image     *Image               `json:"image"`
	Lang      string               `json:"lang"`
	User      *UserByIdModelParent `json:"user"`
	Role      *Role                `json:"role"`
	RoleName  string               `json:"role_name"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}
