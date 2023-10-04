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

package notify

import (
	"context"

	"github.com/e154/smart-home/adaptors"
	m "github.com/e154/smart-home/models"
)

// TemplateBind ...
type TemplateBind struct {
	adaptor *adaptors.Adaptors
}

// NewTemplateBind ...
func NewTemplateBind(adaptor *adaptors.Adaptors) *TemplateBind {
	return &TemplateBind{
		adaptor: adaptor,
	}
}

// Render ...
func (t *TemplateBind) Render(templateName string, params map[string]interface{}) *m.TemplateRender {
	render, err := t.adaptor.Template.Render(context.Background(), templateName, params)
	if err != nil {
		return nil
	}
	return render
}
