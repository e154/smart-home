// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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
	"fmt"

	"github.com/e154/smart-home/internal/common"
	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.TemplateRepo = (*Template)(nil)

// Template ...
type Template struct {
	table *db.Templates
	db    *gorm.DB
}

// GetTemplateAdaptor ...
func GetTemplateAdaptor(d *gorm.DB) *Template {
	return &Template{
		table: &db.Templates{&db.Common{Db: d}},
		db:    d,
	}
}

// UpdateOrCreate ...
func (n *Template) UpdateOrCreate(ctx context.Context, ver *models.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.UpdateOrCreate(ctx, dbVer); err != nil {
		return
	}

	return
}

// Create ...
func (n *Template) Create(ctx context.Context, ver *models.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.Create(ctx, dbVer); err != nil {
		return
	}

	return
}

// UpdateStatus ...
func (n *Template) UpdateStatus(ctx context.Context, ver *models.Template) (err error) {

	dbVer := n.toDb(ver)
	if err = n.table.UpdateStatus(ctx, dbVer); err != nil {
		return
	}

	return
}

// GetList ...
func (n *Template) GetList(ctx context.Context, templateType models.TemplateType) (items []*models.Template, err error) {

	var dbItems []*db.Template
	if dbItems, err = n.table.GetList(ctx, templateType.String()); err != nil {
		return
	}

	items = make([]*models.Template, 0, len(dbItems))
	for _, dbVer := range dbItems {
		items = append(items, n.fromDb(dbVer))
	}

	return
}

// GetByName ...
func (n *Template) GetByName(ctx context.Context, name string) (ver *models.Template, err error) {

	var dbVer *db.Template
	if dbVer, err = n.table.GetByName(ctx, name, "template"); err != nil {
		return
	}

	ver = n.fromDb(dbVer)
	return
}

// GetItemByName ...
func (n *Template) GetItemByName(ctx context.Context, name string) (ver *models.Template, err error) {

	var dbVer *db.Template
	if dbVer, err = n.table.GetByName(ctx, name, "item"); err != nil {
		return
	}

	ver = n.fromDb(dbVer)
	return
}

// GetItemsSortedList ...
func (n *Template) GetItemsSortedList(ctx context.Context) (count int64, items []string, err error) {
	count, items, err = n.table.GetItemsSortedList(ctx)
	return
}

// Delete ...
func (n *Template) Delete(ctx context.Context, name string) (err error) {
	err = n.table.Delete(ctx, name)
	return
}

// GetItemsTree ...
func (n *Template) GetItemsTree(ctx context.Context) (tree []*models.TemplateTree, err error) {

	var dbTree []*db.TemplateTree
	if dbTree, err = n.table.GetItemsTree(ctx); err != nil {
		return
	}

	tree = make([]*models.TemplateTree, 0, len(dbTree))
	err = common.Copy(&tree, &dbTree, common.JsonEngine)

	return
}

// UpdateItemsTree ...
func (n *Template) UpdateItemsTree(ctx context.Context, tree []*models.TemplateTree) (err error) {

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
func (n *Template) Search(ctx context.Context, query string, limit, offset int) (list []*models.Template, total int64, err error) {
	var dbList []*db.Template
	if dbList, total, err = n.table.Search(ctx, query, limit, offset); err != nil {
		return
	}

	list = make([]*models.Template, 0)
	for _, dbVer := range dbList {
		ver := n.fromDb(dbVer)
		_ = n.GetMarkers(ctx, ver)
		list = append(list, ver)
	}

	return
}

// GetMarkers ...
func (n *Template) GetMarkers(ctx context.Context, template *models.Template) (err error) {

	var templateContent *models.TemplateContent
	var items models.Templates

	if templateContent, err = template.GetTemplate(); err != nil {
		return
	}

	if items, err = n.GetList(ctx, models.TemplateTypeItem); err != nil {
		return
	}

	if _, e := template.GetMarkers(items, templateContent); e != nil {
		err = fmt.Errorf("%s: get markers", e.Error())
	}

	return
}

// Render ...
func (n *Template) Render(ctx context.Context, name string, params map[string]interface{}) (render *models.TemplateRender, err error) {

	var item *models.Template
	var template *models.TemplateContent
	var items models.Templates

	if item, err = n.GetByName(ctx, name); err != nil {
		return
	}

	if template, err = item.GetTemplate(); err != nil {
		return
	}

	if items, err = n.GetList(ctx, models.TemplateTypeItem); err != nil {
		return
	}

	render, err = models.RenderTemplate(items, template, params)

	return
}

func (n *Template) fromDb(dbVer *db.Template) (ver *models.Template) {
	ver = &models.Template{
		Name:        dbVer.Name,
		Description: dbVer.Description,
		Content:     dbVer.Content,
		Status:      models.TemplateStatus(dbVer.Status),
		Type:        models.TemplateType(dbVer.Type),
		ParentName:  dbVer.ParentName,
		CreatedAt:   dbVer.CreatedAt,
		UpdatedAt:   dbVer.UpdatedAt,
	}
	return
}

func (n *Template) toDb(ver *models.Template) (dbVer *db.Template) {
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
