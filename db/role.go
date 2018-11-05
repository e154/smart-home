package db

import "time"

type Role struct {
	Name        string `gorm:"primary_key"`
	Description string
	Parent      *Role
	Children    []*Role
	Permissions []*Permission
	CreatedAt   time.Time
	UpdateAt    *time.Time
	//AccessList  map[string][]string `gorm:"-"`
}
