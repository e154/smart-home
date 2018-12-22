package models

type Permission struct {
	Id          int64  `json:"id"`
	RoleName    string `json:"role_name"`
	PackageName string `json:"package_name"`
	LevelName   string `json:"level_name"`
}
