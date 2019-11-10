package models

import "time"

// swagger:model
type NewTemplate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type UpdateTemplate struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	Type        string  `json:"type"`
	ParentName  *string `json:"parent"`
}

// swagger:model
type Template struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	ParentName  *string   `json:"parent"`
	Markers     []string  `json:"markers"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TemplateField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// swagger:model
type TemplateContent struct {
	Items  []string         `json:"items"`
	Title  string           `json:"title"`
	Fields []*TemplateField `json:"fields"`
}
