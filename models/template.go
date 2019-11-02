package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Template struct {
	Validity
	Name         string         `json:"name" valid:"Required;MaxSize(64)" `
	Description  string         `json:"description"`
	Content      string         `json:"content"`
	Status       TemplateStatus `json:"status" valid:"Required;MaxSize(64)"`
	Type         TemplateType   `json:"type" valid:"Required;MaxSize(64)"`
	ParentsCount int            `json:"parents_count"`
	ParentName   *string        `json:"parent"`
	Markers      []string       `json:"markers"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Templates []*Template

func (i Templates) Len() int {
	return len(i)
}

func (i Templates) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i Templates) Less(a, b int) bool {
	return i[a].ParentsCount < i[b].ParentsCount
}

type TemplateField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TemplateContent struct {
	Items  []string         `json:"items"`
	Title  string           `json:"title"`
	Fields []*TemplateField `json:"fields"`
}

func TemplateGetParents(items Templates, result *Templates, s string) {

	for _, item := range items {
		if item.Name == s {
			var exist bool
			for _, v := range *result {
				if v.Name == item.Name {
					exist = true
				}
			}
			if !exist {
				*result = append(*result, item)
			}
			var parent string
			if item.ParentName != nil {
				parent = *item.ParentName
			}
			TemplateGetParents(items, result, parent)
		}
	}
}

func (i *Template) GetMarkers() (markers []string, err error) {

	if i.Type != TemplateTypeTemplate {
		return
	}

	tpl := &TemplateContent{}
	if err = json.Unmarshal([]byte(i.Content), tpl); err != nil {
		fmt.Println(err.Error())
		return
	}

	reg, _ := regexp.CompilePOSIX(`\[{1}([a-zA-Z\-_0-9:]*)\]{1}`)

	var findMarkers func(string)

	findMarkers = func(s string) {

		ms := reg.FindAllStringSubmatch(s, -1)
		for _, m := range ms {
			if strings.Contains(m[1], "content") || strings.Contains(m[1], "block") {
				continue
			}

			var exist bool
			for _, _m := range markers {
				if _m == m[1] {
					exist = true
				}
			}

			if !exist {
				markers = append(markers, m[1])
			}
		}
	}

	for _, field := range tpl.Fields {
		findMarkers(field.Value)
	}

	findMarkers(tpl.Title)

	i.Markers = markers

	return
}

func (i *Template) GetTemplate() (tpl *TemplateContent, err error) {

	tpl = new(TemplateContent)
	err = json.Unmarshal([]byte(i.Content), tpl)
	return
}

type TemplateRender struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
