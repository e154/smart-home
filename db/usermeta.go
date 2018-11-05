package db

type UserMeta struct {
	Id    int64 `gorm:"primary_key"`
	User  *User
	Key   string
	Value string
}
