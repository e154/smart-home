package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

type Templates struct {
	Db *gorm.DB
}

type TemplateTree struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Nodes       []*TemplateTree `json:"nodes"`
}

type Template struct {
	Name        string `gorm:"primary_key"`
	Description string
	Content     string
	Status      string
	Type        string
	ParentName  *string `gorm:"column:parent"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (d *Template) TableName() string {
	return "templates"
}

func (n Templates) UpdateOrCreate(tpl *Template) error {

	if err := n.Update(tpl); err != nil {
		if err = n.Db.Create(tpl).Error; err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func (n Templates) Create(tpl *Template) error {

	if err := n.Db.Create(tpl).Error; err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (n Templates) GetByName(name, itemType string) (*Template, error) {

	tpl := &Template{}
	err := n.Db.Model(tpl).
		Where("name = ? and type = ?", name, itemType).
		First(&tpl).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return tpl, nil
}

func (n Templates) GetItemsSortedList() (count int64, newItems []string, err error) {

	items := make([]*Template, 0)
	err = n.Db.Model(&Template{}).
		Where("type = 'item' and status = 'active'").
		Find(&items).
		Error

	if err != nil {
		err = errors.New(err.Error())
		return
	}

	newItems = make([]string, 0)

	treeGetEndPoints := func(i []*Template, t *[]string) {
		for _, v := range i {
			var exist bool
			for _, k := range i {
				if k.ParentName == &v.Name {
					exist = true
					break
				}
			}

			if !exist {
				*t = append(*t, v.Name)
			}
		}
	};
	treeGetEndPoints(items, &newItems)
	count = int64(len(newItems))

	return
}

func (n Templates) Update(m *Template) error {
	err := n.Db.Model(&Template{Name: m.Name}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"type":        m.Type,
		"content":     m.Content,
		"parent":      m.ParentName,
	}).Error

	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (n Templates) UpdateStatus(m *Template) error {
	err := n.Db.Model(&Template{Name: m.Name}).Updates(map[string]interface{}{
		"status": m.Status,
	}).Error

	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (n Templates) Delete(name string) (err error) {
	err = n.Db.Delete(&Template{Name: name}).Error
	return
}

func (n Templates) GetItemsTree() (tree *TemplateTree, err error) {

	var items []*Template
	if items, err = n.GetList("item"); err != nil {
		return
	}

	tree = &TemplateTree{}
	for _, item := range items {
		if item.ParentName == nil {
			tree.Description = item.Description
			tree.Name = item.Name
		}
	}

	n.renderTreeRecursive(items, tree, tree.Name)

	return
}

func (n Templates) GetList(templateType string) ([]*Template, error) {
	items := make([]*Template, 0)
	err := n.Db.Model(&Template{}).
		Where("type = ?", templateType).
		Find(&items).
		Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return items, nil
}

func (n *Templates) Search(query string, limit, offset int) (items []*Template, total int64, err error) {

	q := n.Db.Model(&Template{}).
		Where("name LIKE ?", "%"+query+"%").
		Order("name ASC")

	if err = q.Count(&total).Error; err != nil {
		err = errors.New(err.Error())
		return
	}

	items = make([]*Template, 0)
	if err = q.Find(&items).Error; err != nil {
		err = errors.New(err.Error())
	}

	return
}

func (n Templates) renderTreeRecursive(i []*Template, t *TemplateTree, c string) {

	for _, item := range i {
		if item.ParentName != nil && *item.ParentName == c {
			tree := &TemplateTree{}
			tree.Name = item.Name
			tree.Description = item.Description
			tree.Nodes = make([]*TemplateTree, 0) // fix - nodes: null
			tree.Status = item.Status
			t.Nodes = append(t.Nodes, tree)
			n.renderTreeRecursive(i, tree, item.Name)
		}
	}

	return
}

func (n Templates) UpdateItemsTree(tree []*TemplateTree, parent string) error {

	for _, v := range tree {
		if parent != "" {
			go n.emailItemParentUpdate(v.Name, parent)
		}
		n.updateTreeRecursive(v.Nodes, v.Name)
	}

	return nil
}

func (n Templates) emailItemParentUpdate(name, parent string) error {

	err := n.Db.Model(&Template{}).
		Where("name = ?", name).
		Updates(map[string]interface{}{
			"parent": parent,
		}).Error

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (n Templates) updateTreeRecursive(t []*TemplateTree, parent string) {

	for _, v := range t {
		if parent != "" {
			go n.emailItemParentUpdate(v.Name, parent)
		}
		n.updateTreeRecursive(v.Nodes, v.Name)
	}

}
