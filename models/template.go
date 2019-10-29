package models

import (
	"time"
)

type Template struct {
	Validity
	Name        string         `json:"name" valid:"Required;MaxSize(64)" `
	Description string         `json:"description"`
	Content     string         `json:"content"`
	Status      TemplateStatus `json:"status" valid:"Required;MaxSize(64)"`
	Type        TemplateType   `json:"type" valid:"Required;MaxSize(64)"`
	ParentName  *string        `json:"parent"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
