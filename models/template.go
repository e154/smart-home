package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
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

// TODO add recursive search
func (i *Template) GetMarkers() (markers []string, err error) {

	markers = make([]string, 0)

	if i.Type != TemplateTypeTemplate {
		return
	}

	tpl := &TemplateContent{}
	if err = json.Unmarshal([]byte(i.Content), tpl); err != nil {
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

func RenderTemplate(items []*Template, template *TemplateContent, params map[string]interface{}) (render *TemplateRender, err error) {

	parents := Templates{}
	for _, item := range template.Items {
		TemplateGetParents(items, &parents, item)
	}

	for _, item := range parents {
		p := make(Templates, 0)
		TemplateGetParents(parents, &p, item.Name)
		item.ParentsCount = len(p)
	}

	sort.Sort(parents)

	var buf string

	// замена [xxxx:block] на реальные блоки
	for key, item := range parents {
		if item.Status != "active" {
			continue
		}

		if key == 0 {
			buf = item.Content
		} else {
			buf = strings.Replace(buf, fmt.Sprintf("[%s:block]", item.Name), item.Content, -1)
		}
	}

	// поиск маркера [xxx:content] и замена на контейнер с возможностью редактирования
	reg := regexp.MustCompile(`(\[{1}[a-zA-Z0-9_\-]{2,64}\:content\]{1})`)
	reg2 := regexp.MustCompile(`(\[{1})([a-zA-Z0-9_\-]{2,64})(\:)([content]+)([\]]{1})`)
	markers := reg.FindAllString(buf, -1)
	for _, m := range markers {
		marker := reg2.FindStringSubmatch(m)[2]

		f := m
		for _, field := range template.Fields {
			if field.Name == marker {
				if utf8.RuneCountInString(field.Value) > 0 {
					f = field.Value
				}
			}
		}

		buf = strings.Replace(buf, m, f, -1)
	}

	// скрыть не использованные маркеры [xxxx:block] блоков
	reg = regexp.MustCompile(`(\[{1}[a-zA-Z0-9_\-]{2,64}\:block\]{1})`)
	blocks := reg.FindAllString(buf, -1)
	for _, block := range blocks {
		buf = strings.Replace(buf, block, "", -1)
	}

	// заполнение формы
	title := template.Title
	if params != nil {
		for k, v := range params {
			buf = strings.Replace(buf, fmt.Sprintf("[%s]", k), fmt.Sprintf("%v", v), -1)
			title = strings.Replace(title, fmt.Sprintf("[%s]", k), fmt.Sprintf("%v", v), -1)
		}
	}

	render = &TemplateRender{
		Subject: title,
		Body:    buf,
	}

	return
}

func PreviewTemplate(items []*Template, template *TemplateContent) (data string, err error) {

	parents := Templates{}
	for _, item := range template.Items {
		TemplateGetParents(items, &parents, item)
	}

	for _, item := range parents {
		p := make(Templates, 0)
		TemplateGetParents(parents, &p, item.Name)
		item.ParentsCount = len(p)
	}

	sort.Sort(parents)

	// замена [xxxx:block] на реальные блоки
	for key, item := range parents {
		if item.Status != "active" {
			continue
		}

		if key == 0 {
			data = item.Content
		} else {
			data = strings.Replace(data, fmt.Sprintf("[%s:block]", item.Name), item.Content, -1)
		}
	}

	// поиск маркера [xxx:content] и замена на контейнер с возможностью редактирования
	reg := regexp.MustCompile(`(\[{1}[a-zA-Z0-9_\-]{2,64}\:content\]{1})`)
	reg2 := regexp.MustCompile(`(\[{1})([a-zA-Z0-9_\-]{2,64})(\:)([content]+)([\]]{1})`)
	markers := reg.FindAllString(data, -1)
	for _, m := range markers {
		marker := reg2.FindStringSubmatch(m)[2]

		f := m
		for _, field := range template.Fields {
			if field.Name == marker {
				if utf8.RuneCountInString(field.Value) > 0 {
					f = field.Value
				}
			}
		}

		data = strings.Replace(data, m, fmt.Sprintf("<div class=\"edit_inline\" data-id=\"%s\">%s</div>", marker, f), -1)
	}

	// скрыть не использованные маркеры [xxxx:block] блоков
	reg = regexp.MustCompile(`(\[{1}[a-zA-Z0-9_\-]{2,64}\:block\]{1})`)
	blocks := reg.FindAllString(data, -1)
	for _, block := range blocks {
		data = strings.Replace(data, block, "", -1)
	}

	return
}