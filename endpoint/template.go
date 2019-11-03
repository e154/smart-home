package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"github.com/pkg/errors"
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

func (t *TemplateEndpoint) UpdateStatus(params *m.Template) (errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	if err = t.adaptors.Template.UpdateStatus(params); err != nil {
		return
	}
	return
}

func (t *TemplateEndpoint) GetByName(name string) (result *m.Template, err error) {
	if result, err = t.adaptors.Template.GetByName(name); err != nil {
		return
	}

	if _, e := result.GetMarkers(); e != nil {
		err = errors.Wrap(e, "get markers")
	}
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

func (t *TemplateEndpoint) Preview(template *m.TemplateContent) (data string, err error) {

	var items []*m.Template
	if items, err = t.adaptors.Template.GetList(m.TemplateTypeItem); err != nil {
		return
	}

	data, err = m.PreviewTemplate(items, template)

	return
}
