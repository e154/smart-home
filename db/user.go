package db

import (
	"time"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"database/sql"
)

const HISTORY_MAX = 8

type Users struct {
	Db *gorm.DB
}

type User struct {
	Id                  int64 `gorm:"primary_key"`
	Nickname            string
	FirstName           string
	LastName            string
	EncryptedPassword   string
	Email               string
	Status              string
	ResetPasswordToken  string
	AuthenticationToken string
	Image               *Image
	ImageId             sql.NullInt64
	SignInCount         int64
	CurrentSignInIp     string
	LastSignInIp        string
	Lang                string
	User                *User
	UserId              sql.NullInt64
	Role                *Role
	RoleId              int64
	Meta                []*UserMeta
	ResetPasswordSentAt time.Time
	CurrentSignInAt     time.Time
	LastSignInAt        time.Time
	CreatedAt           time.Time
	UpdatedAt           *time.Time
	DeletedAt           *time.Time
	History             json.RawMessage `gorm:"type:jsonb;not null"`
}

func (m *User) TableName() string {
	return "users"
}
