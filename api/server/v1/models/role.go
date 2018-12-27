package models

import "time"

type NewRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	}
}

type UpdateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	}
}

type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	}
	Children   []*Role             `json:"children"`
	AccessList map[string][]string `json:"access_list"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type RoleListModel struct {
	Items []UserShotModel
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}

type SearchRoleResponse struct {
	Roles []*Role `json:"roles"`
}

type AccessItem struct {
	Actions     []string `json:"actions"`
	Method      string   `json:"method"`
	Description string   `json:"description"`
	RoleName    string   `json:"role_name"`
}

type AccessLevels map[string]AccessItem
type AccessList map[string]AccessLevels

type AccessListDiff map[string]map[string]bool
