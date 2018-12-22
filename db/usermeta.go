package db

import "github.com/jinzhu/gorm"

type UserMetas struct {
	Db *gorm.DB
}

type UserMeta struct {
	Id     int64 `gorm:"primary_key"`
	User   *User
	UserId int64
	Key    string
	Value  string
}

func (m *UserMeta) TableName() string {
	return "user_metas"
}
