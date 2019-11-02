package models

import "time"

// swagger:model
type NewTemplateItem struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type UpdateTemplateItem struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type UpdateTemplateItemStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// swagger:model
type TemplateItem struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	ParentName  *string   `json:"parent"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// swagger:model
type TemplateTree struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Nodes       []*TemplateTree `json:"nodes"`
}

// swagger:model
type UpdateTemplateTree []*TemplateTree
