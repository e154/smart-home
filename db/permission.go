package db

type Permission struct {
	Id          int64 `gorm:"primary_key"`
	Role        *Role
	PackageName string
	LevelName   string
}
