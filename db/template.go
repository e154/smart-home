// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/e154/smart-home/common/apperr"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Templates ...
type Templates struct {
	Db *gorm.DB
}

// TemplateTree ...
type TemplateTree struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      string          `json:"status"`
	Nodes       []*TemplateTree `json:"nodes"`
}

// Template ...
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

// TableName ...
func (d *Template) TableName() string {
	return "templates"
}

// UpdateOrCreate ...
func (n Templates) UpdateOrCreate(ctx context.Context, tpl *Template) (err error) {

	if err = n.Db.WithContext(ctx).Create(tpl).Error; err != nil {
		if err = n.Update(ctx, tpl); err != nil {
			err = errors.Wrap(apperr.ErrTemplateUpdate, err.Error())
			return
		}
	}

	return
}

// Create ...
func (n Templates) Create(ctx context.Context, tpl *Template) (err error) {
	if err = n.Db.WithContext(ctx).Create(tpl).Error; err != nil {
		err = errors.Wrap(apperr.ErrTemplateAdd, err.Error())
	}
	return
}

// GetByName ...
func (n Templates) GetByName(ctx context.Context, name, itemType string) (*Template, error) {

	tpl := &Template{}
	err := n.Db.WithContext(ctx).Model(tpl).
		Where("name = ? and type = ?", name, itemType).
		First(&tpl).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.Wrap(apperr.ErrTemplateNotFound, fmt.Sprintf("name \"%s\"", name))
			return nil, err
		}
		err = errors.Wrap(apperr.ErrTemplateGet, err.Error())
	}

	return tpl, nil
}

// GetItemsSortedList ...
func (n Templates) GetItemsSortedList(ctx context.Context) (count int64, newItems []string, err error) {

	items := make([]*Template, 0)
	err = n.Db.WithContext(ctx).Model(&Template{}).
		Where("type = 'item' and status = 'active'").
		Find(&items).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrTemplateList, err.Error())
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
	}
	treeGetEndPoints(items, &newItems)
	count = int64(len(newItems))

	return
}

// Update ...
func (n Templates) Update(ctx context.Context, m *Template) error {
	err := n.Db.WithContext(ctx).Model(&Template{Name: m.Name}).Updates(map[string]interface{}{
		"name":        m.Name,
		"description": m.Description,
		"status":      m.Status,
		"type":        m.Type,
		"content":     m.Content,
		"parent":      m.ParentName,
	}).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrTemplateUpdate, err.Error())
	}
	return nil
}

// UpdateStatus ...
func (n Templates) UpdateStatus(ctx context.Context, m *Template) error {
	err := n.Db.WithContext(ctx).Model(&Template{Name: m.Name}).Updates(map[string]interface{}{
		"status": m.Status,
	}).Error

	if err != nil {
		err = errors.Wrap(apperr.ErrTemplateUpdate, err.Error())
	}
	return nil
}

// Delete ...
func (n Templates) Delete(ctx context.Context, name string) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&Template{Name: name}).Error; err != nil {
		err = errors.Wrap(apperr.ErrTemplateDelete, err.Error())
	}
	return
}

// GetItemsTree ...
func (n Templates) GetItemsTree(ctx context.Context) (tree []*TemplateTree, err error) {

	var items []*Template
	if items, err = n.GetList(ctx, "item"); err != nil {
		err = errors.Wrap(apperr.ErrTemplateGet, err.Error())
		return
	}

	tree = make([]*TemplateTree, 0)
	for _, item := range items {
		if item.ParentName == nil {
			branch := &TemplateTree{
				Description: item.Description,
				Name:        item.Name,
				Nodes:       make([]*TemplateTree, 0),
				Status:      item.Status,
			}
			n.renderTreeRecursive(ctx, items, branch, branch.Name)
			tree = append(tree, branch)
		}
	}

	return
}

// GetList ...
func (n Templates) GetList(ctx context.Context, templateType string) ([]*Template, error) {
	items := make([]*Template, 0)
	err := n.Db.WithContext(ctx).Model(&Template{}).
		Where("type = ?", templateType).
		Find(&items).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrTemplateGet, err.Error())
	}

	return items, nil
}

// Search ...
func (n *Templates) Search(ctx context.Context, query string, limit, offset int) (items []*Template, total int64, err error) {

	q := n.Db.WithContext(ctx).Model(&Template{}).
		Where("name LIKE ?", "%"+query+"%").
		Where("type = 'template'")

	if err = q.Count(&total).Error; err != nil {
		err = errors.Wrap(apperr.ErrTemplateSearch, err.Error())
		return
	}

	q = q.
		Limit(limit).
		Offset(offset).
		Order("name ASC")

	items = make([]*Template, 0)
	if err = q.Find(&items).Error; err != nil {
		err = errors.Wrap(apperr.ErrTemplateSearch, err.Error())
	}

	return
}

func (n Templates) renderTreeRecursive(ctx context.Context, i []*Template, t *TemplateTree, c string) {

	for _, item := range i {
		if item.ParentName != nil && *item.ParentName == c {
			tree := &TemplateTree{}
			tree.Name = item.Name
			tree.Description = item.Description
			tree.Nodes = make([]*TemplateTree, 0) // fix - nodes: null
			tree.Status = item.Status
			t.Nodes = append(t.Nodes, tree)
			n.renderTreeRecursive(ctx, i, tree, item.Name)
		}
	}
}

// UpdateItemsTree ...
func (n Templates) UpdateItemsTree(ctx context.Context, tree []*TemplateTree, parent string) error {

	for _, v := range tree {
		if parent != "" {
			go n.emailItemParentUpdate(ctx, v.Name, parent)
		}

		err := n.Db.WithContext(ctx).Model(&Template{Name: v.Name}).Updates(map[string]interface{}{
			"parent": nil,
		}).Error
		if err != nil {
			err = errors.Wrap(apperr.ErrTemplateUpdate, err.Error())
		}

		if len(v.Nodes) == 0 {
			continue
		}

		n.updateTreeRecursive(ctx, v.Nodes, v.Name)
	}

	return nil
}

func (n Templates) emailItemParentUpdate(ctx context.Context, name, parent string) {

	_ = n.Db.WithContext(ctx).Model(&Template{}).
		Where("name = ?", name).
		Updates(map[string]interface{}{
			"parent": parent,
		}).Error
}

func (n Templates) updateTreeRecursive(ctx context.Context, t []*TemplateTree, parent string) {

	for _, v := range t {
		if parent != "" {
			go n.emailItemParentUpdate(ctx, v.Name, parent)
		}
		n.updateTreeRecursive(ctx, v.Nodes, v.Name)
	}

}
