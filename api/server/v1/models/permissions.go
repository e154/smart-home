package models

type Permission struct {
	Id          int64  `json:"id"`
	PackageName string `json:"package_name"`
	LevelName   string `json:"level_name"`
}

type Permissions []*Permission