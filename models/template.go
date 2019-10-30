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

type Templates []*Template

type TemplateField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TemplatePreview struct {
	Items  []string         `json:"items"`
	Title  string           `json:"title"`
	Fields []*TemplateField `json:"fields"`
}

func TemplateGetParents(items Templates, result *Templates, s string) {

	for _, item := range items {
		if item.ParentName != nil && *item.ParentName == s {
			var exist bool
			for _, v := range *result {
				if v.Name == item.Name {
					exist = true
				}
			}
			if !exist {
				*result = append(*result, item)
			}
			TemplateGetParents(items, result, *item.ParentName)
		}
	}
}
