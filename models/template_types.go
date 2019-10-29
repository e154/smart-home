package models

type TemplateTree struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Nodes       []*TemplateTree `json:"nodes"`
}

type TemplateStatus string
type TemplateType string

func (s TemplateStatus) String() string {
	return string(s)
}

func (t TemplateType) String() string {
	return string(t)
}

const (
	TemplateStatusActive   = TemplateStatus("active")
	TemplateStatusUnactive = TemplateStatus("unactive")
	TemplateTypeItem       = TemplateType("item")
	TemplateTypeTemplate   = TemplateType("template")
)
