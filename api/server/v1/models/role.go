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

type ResponseRole struct {
	Code ResponseType `json:"code"`
	Data struct {
		Role *Role `json:"role"`
	} `json:"data"`
}

type ResponseRoleList struct {
	Code ResponseType `json:"code"`
	Data struct {
		Items  []*Role `json:"items"`
		Limit  int64   `json:"limit"`
		Offset int64   `json:"offset"`
		Total  int64   `json:"total"`
	} `json:"data"`
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
