package models

type UserMeta struct {
	Id     int64  `json:"id"`
	User   *User  `json:"user"`
	UserId int64  `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
