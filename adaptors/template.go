// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
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

package adaptors

import (
	"context"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ITemplate ...
type ITemplate interface {
	UpdateOrCreate(ctx context.Context, ver *m.Template) (err error)
	Create(ctx context.Context, ver *m.Template) (err error)
	UpdateStatus(ctx context.Context, ver *m.Template) (err error)
	GetList(ctx context.Context, templateType m.TemplateType) (items []*m.Template, err error)
	GetByName(ctx context.Context, name string) (ver *m.Template, err error)
	GetItemByName(ctx context.Context, name string) (ver *m.Template, err error)
	GetItemsSortedList(ctx context.Context, ) (count int64, items []string, err error)
	Delete(ctx context.Context, name string) (err error)
	GetItemsTree(ctx context.Context, ) (tree []*m.TemplateTree, err error)
	UpdateItemsTree(ctx context.Context, tree []*m.TemplateTree) (err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*m.Template, total int64, err error)
	GetMarkers(ctx context.Context, template *m.Template) (err error)
	Render(ctx context.Context, name string, params map[string]interface{}) (render *m.TemplateRender, err error)
	fromDb(dbVer *db.Template) (ver *m.Template)
	toDb(ver *m.Template) (dbVer *db.Template)
}

// Template ...
type Template struct {
	ITemplate
	table *db.Templates
	db    *gorm.DB
}

// GetTemplateAdaptor ...
func GetTemplateAdaptor(d *gorm.DB) ITemplate {
	return &Template{
		table: &db.Templates{Db: d},
		db:    d,
	}
}

// UpdateOrCreate ...
func (n *Template) UpdateOrCreate(ctx context.Context, ver *m.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.UpdateOrCreate(ctx, dbVer); err != nil {
		return
	}

	return
}

// Create ...
func (n *Template) Create(ctx context.Context, ver *m.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.Create(ctx, dbVer); err != nil {
		return
	}

	return
}

// UpdateStatus ...
func (n *Template) UpdateStatus(ctx context.Context, ver *m.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.UpdateStatus(ctx, dbVer); err != nil {
		return
	}

	return
}

// GetList ...
func (n *Template) GetList(ctx context.Context, templateType m.TemplateType) (items []*m.Template, err error) {

	var dbItems []*db.Template
	if dbItems, err = n.table.GetList(ctx, templateType.String()); err != nil {
		return
	}

	items = make([]*m.Template, 0, len(dbItems))
	for _, dbVer := range dbItems {
		items = append(items, n.fromDb(dbVer))
	}

	return
}

// GetByName ...
func (n *Template) GetByName(ctx context.Context, name string) (ver *m.Template, err error) {

	var dbVer *db.Template
	if dbVer, err = n.table.GetByName(ctx, name, "template"); err != nil {
		return
	}

	ver = n.fromDb(dbVer)
	return
}

// GetItemByName ...
func (n *Template) GetItemByName(ctx context.Context, name string) (ver *m.Template, err error) {

	var dbVer *db.Template
	if dbVer, err = n.table.GetByName(ctx, name, "item"); err != nil {
		return
	}

	ver = n.fromDb(dbVer)
	return
}

// GetItemsSortedList ...
func (n *Template) GetItemsSortedList(ctx context.Context, ) (count int64, items []string, err error) {
	count, items, err = n.table.GetItemsSortedList(ctx, )
	return
}

// Delete ...
func (n *Template) Delete(ctx context.Context, name string) (err error) {
	err = n.table.Delete(ctx, name)
	return
}

// GetItemsTree ...
func (n *Template) GetItemsTree(ctx context.Context, ) (tree []*m.TemplateTree, err error) {

	var dbTree []*db.TemplateTree
	if dbTree, err = n.table.GetItemsTree(ctx, ); err != nil {
		return
	}

	tree = make([]*m.TemplateTree, 0, len(dbTree))
	err = common.Copy(&tree, &dbTree, common.JsonEngine)

	return
}

// UpdateItemsTree ...
func (n *Template) UpdateItemsTree(ctx context.Context, tree []*m.TemplateTree) (err error) {

	dbTree := make([]*db.TemplateTree, 0)
	if err = common.Copy(&dbTree, &tree, common.JsonEngine); err != nil {
		return
	}

	if err = n.table.UpdateItemsTree(ctx, dbTree, ""); err != nil {
		return
	}

	return
}

// Search ...
func (n *Template) Search(ctx context.Context, query string, limit, offset int) (list []*m.Template, total int64, err error) {
	var dbList []*db.Template
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]*m.Template, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		_ = n.GetMarkers(ctx, ver)
		list = append(list, ver)
	}

	return
}

// GetMarkers ...
func (n *Template) GetMarkers(ctx context.Context, template *m.Template) (err error) {

	var templateContent *m.TemplateContent
	var items m.Templates

	if templateContent, err = template.GetTemplate(); err != nil {
		return
	}

	if items, err = n.GetList(ctx, m.TemplateTypeItem); err != nil {
		return
	}

	if _, e := template.GetMarkers(items, templateContent); e != nil {
		err = errors.Wrap(e, "get markers")
	}

	return
}

// Render ...
func (n *Template) Render(ctx context.Context, name string, params map[string]interface{}) (render *m.TemplateRender, err error) {

	var item *m.Template
	var template *m.TemplateContent
	var items m.Templates

	if item, err = n.GetByName(ctx, name); err != nil {
		return
	}

	if template, err = item.GetTemplate(); err != nil {
		return
	}

	if items, err = n.GetList(ctx, m.TemplateTypeItem); err != nil {
		return
	}

	render, err = m.RenderTemplate(items, template, params)

	return
}

func (n *Template) fromDb(dbVer *db.Template) (ver *m.Template) {
	ver = &m.Template{
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Content:     dbVer.Content,
		Status:      m.TemplateStatus(dbVer.Status),
		Type:        m.TemplateType(dbVer.Type),
		ParentName:  dbVer.ParentName,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	return
}

func (n *Template) toDb(ver *m.Template) (dbVer *db.Template) {
	dbVer = &db.Template{
		Name:        ver.Name,
		Description: ver.Description,
		Content:     ver.Content,
		Status:      ver.Status.String(),
		Type:        ver.Type.String(),
		ParentName:  ver.ParentName,
		CreatedAt:   ver.CreatedAt,
		UpdatedAt:   ver.UpdatedAt,
	}
	return
}
