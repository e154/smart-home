package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
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

func (t *TemplateEndpoint) GetItemsSortedList() (count int64, items []string, err error) {
	count, items, err = t.adaptors.Template.GetItemsSortedList()
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
