package endpoint

import (
	"fmt"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"regexp"
	"strings"
	"unicode/utf8"
)

type TemplateEndpoint struct {
	*CommonEndpoint
}

func NewTemplateEndpoint(common *CommonEndpoint) *TemplateEndpoint {
	return &TemplateEndpoint{
		CommonEndpoint: common,
	}
}

func (t *TemplateEndpoint) UpdateOrCreate(params *m.Template) (errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = t.adaptors.Template.UpdateOrCreate(params); err != nil {
		return
	}
	return
}

func (t *TemplateEndpoint) GetByName(name string) (result *m.Template, err error) {
	result, err = t.adaptors.Template.GetByName(name)
	return
}

func (t *TemplateEndpoint) GetItemByName(name string) (result *m.Template, err error) {
	result, err = t.adaptors.Template.GetItemByName(name)
	return
}

func (t *TemplateEndpoint) GetItemsSortedList() (count int64, items []string, err error) {
	count, items, err = t.adaptors.Template.GetItemsSortedList()
	return
}

func (t *TemplateEndpoint) GetList() (count int64, templates []*m.Template, err error) {
	templates, err = t.adaptors.Template.GetList(m.TemplateTypeTemplate)
	return
}

func (t *TemplateEndpoint) Delete(name string) (err error) {
	err = t.adaptors.Template.Delete(name)
	return
}

func (t *TemplateEndpoint) GetItemsTree() (tree *m.TemplateTree, err error) {
	tree, err = t.adaptors.Template.GetItemsTree()
	return
}

func (t *TemplateEndpoint) UpdateItemsTree(tree []*m.TemplateTree) (err error) {
	err = t.adaptors.Template.UpdateItemsTree(tree)
	return
}

func (t *TemplateEndpoint) Search(query string, limit, offset int) (result []*m.Template, total int64, err error) {
	result, total, err = t.adaptors.Template.Search(query, limit, offset)
	return
}

func (t *TemplateEndpoint) Preview(templatePreview *m.TemplatePreview) (data string, err error) {

	var template *m.TemplatePreview

	var items []*m.Template
	if items, err = t.adaptors.Template.GetList(m.TemplateTypeItem); err != nil {
		return
	}

	result := m.Templates{}
	for _, item := range templatePreview.Items {
		m.TemplateGetParents(items, &result, item)
	}

	//sort.Sort(result)

	// замена [xxxx:block] на реальные блоки
	for key, item := range result {
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
	reg := regexp.MustCompile(`(\[{1}[a-z]{2,64}\:content\]{1})`)
	reg2 := regexp.MustCompile(`(\[{1})([a-z]{2,64})(\:)([content]+)([\]]{1})`)
	markers := reg.FindAllString(data, -1)
	for _, m := range markers {
		marker := reg2.FindStringSubmatch(m)[2]

		var f string = m
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
	reg = regexp.MustCompile(`(\[{1}[a-z]{2,64}\:block\]{1})`)
	blocks := reg.FindAllString(data, -1)
	for _, block := range blocks {
		data = strings.Replace(data, block, "", -1)
	}

	return
}
