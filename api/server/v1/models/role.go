package models

import "time"

// swagger:model
type NewRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	} `json:"parent"`
}

// swagger:model
type UpdateRole struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	} `json:"parent"`
}

// swagger:model
type Role struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parent      struct {
		Name string `json:"name"`
	} `json:"parent"`
	Children   []*Role             `json:"children"`
	AccessList map[string][]string `json:"access_list"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type AccessItem struct {
	Actions     []string `json:"actions"`
	Method      string   `json:"method"`
	Description string   `json:"description"`
	RoleName    string   `json:"role_name"`
}

type AccessLevels map[string]AccessItem

// swagger:model
type AccessList map[string]AccessLevels

// swagger:model
type AccessListDiff map[string]map[string]bool
