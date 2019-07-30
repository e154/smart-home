package models

import (
	"database/sql"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	"time"
)

const HistoryMax = 8

type User struct {
	Id                  int64          `json:"id"`
	Nickname            string         `json:"nickname" valid:"Required;MinSize(3);MaxSize(255)"`
	FirstName           string         `json:"first_name" valid:"MaxSize(255)"`
	LastName            string         `json:"last_name" valid:"MaxSize(255)"`
	EncryptedPassword   string         `json:"encrypted_password,omitempty"`
	Email               string         `json:"email" valid:"Required;Email"`
	Status              string         `json:"status" valid:"MaxSize(255)"`
	ResetPasswordToken  string         `json:"-,omitempty"`
	AuthenticationToken *string        `json:"authentication_token,omitempty"`
	Image               *Image         `json:"image,omitempty"`
	ImageId             sql.NullInt64  `json:"image_id"`
	SignInCount         int64          `json:"sign_in_count,omitempty"`
	CurrentSignInIp     string         `json:"current_sign_in_ip,omitempty"`
	LastSignInIp        string         `json:"last_sign_in_ip,omitempty"`
	Lang                string         `json:"lang"`
	User                *User          `json:"user,omitempty"`
	UserId              sql.NullInt64  `json:"user_id"`
	Role                *Role          `json:"role"`
	RoleName            string         `json:"role_name,omitempty"`
	Meta                []*UserMeta    `json:"meta,omitempty"`
	ResetPasswordSentAt *time.Time     `json:"reset_password_sent_at,omitempty"`
	CurrentSignInAt     *time.Time     `json:"current_sign_in_at,omitempty"`
	LastSignInAt        *time.Time     `json:"last_sign_in_at,omitempty"`
	CreatedAt           time.Time      `json:"created_at,omitempty"`
	UpdatedAt           time.Time      `json:"updated_at,omitempty"`
	DeletedAt           *time.Time     `json:"deleted_at"`
	History             []*UserHistory `json:"history,omitempty"`
}

func (u *User) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(u); !ok {
		errs = valid.Errors
	}

	return
}

func (u *User) UpdateHistory(t time.Time, ipv4 string) {

	l := len(u.History)
	if l > HistoryMax {
		u.History = u.History[l-HistoryMax:]
	}

	u.History = append(u.History, &UserHistory{Ip: ipv4, Time: t})
}

func (u User) CheckPass(password string) (ok bool) {
	ok = common.CheckPasswordHash(password, u.EncryptedPassword)
	return
}

func (u *User) SetPass(password string) (err error) {
	u.EncryptedPassword, err = common.HashPassword(password)
	return
}
