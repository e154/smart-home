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

	"github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"

	"gorm.io/gorm"
)

var _ adaptors.ScriptVersionRepo = (*ScriptVersion)(nil)

// ScriptVersion ...
type ScriptVersion struct {
	table *db.ScriptVersions
	db    *gorm.DB
}

// GetScriptVersionAdaptor ...
func GetScriptVersionAdaptor(d *gorm.DB) *ScriptVersion {
	return &ScriptVersion{
		table: &db.ScriptVersions{&db.Common{Db: d}},
		db:    d,
	}
}

// Delete ...
func (n *ScriptVersion) Delete(ctx context.Context, id int64) (err error) {
	err = n.table.Delete(ctx, id)
	return
}
