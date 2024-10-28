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

	"github.com/e154/smart-home/pkg/models"
)

// TemplateRepo ...
type TemplateRepo interface {
	UpdateOrCreate(ctx context.Context, ver *models.Template) (err error)
	Create(ctx context.Context, ver *models.Template) (err error)
	UpdateStatus(ctx context.Context, ver *models.Template) (err error)
	GetList(ctx context.Context, templateType models.TemplateType) (items []*models.Template, err error)
	GetByName(ctx context.Context, name string) (ver *models.Template, err error)
	GetItemByName(ctx context.Context, name string) (ver *models.Template, err error)
	GetItemsSortedList(ctx context.Context) (count int64, items []string, err error)
	Delete(ctx context.Context, name string) (err error)
	GetItemsTree(ctx context.Context) (tree []*models.TemplateTree, err error)
	UpdateItemsTree(ctx context.Context, tree []*models.TemplateTree) (err error)
	Search(ctx context.Context, query string, limit, offset int) (list []*models.Template, total int64, err error)
	GetMarkers(ctx context.Context, template *models.Template) (err error)
	Render(ctx context.Context, name string, params map[string]interface{}) (render *models.TemplateRender, err error)
}
