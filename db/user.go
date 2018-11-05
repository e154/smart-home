package db

import "time"

type User struct {
	Id                  int64 `gorm:"primary_key"`
	Nickname            string
	FirstName           string
	LastName            string
	EncryptedPassword   string
	Email               string
	HistoryStr          string
	Status              string
	ResetPasswordToken  string
	AuthenticationToken string
	Avatar              *Image
	SignInCount         int64
	CurrentSignInIp     string
	LastSignInIp        string
	Lang                string
	CreatedBy           *User
	Role                *Role
	Meta                []*UserMeta
	ResetPasswordSentAt time.Time
	CurrentSignInAt     time.Time
	LastSignInAt        time.Time
	CreatedAt           time.Time
	UpdateAt            *time.Time
	Deleted             *time.Time
	//History             []*UserHistory `gorm:"-"`
}

const HISTORY_MAX = 8
