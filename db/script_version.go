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
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
)

// ScriptVersions ...
type ScriptVersions struct {
	Db *gorm.DB
}

// ScriptVersion ...
type ScriptVersion struct {
	Id        int64 `gorm:"primary_key"`
	Lang      ScriptLang
	Source    string
	ScriptId  int64
	Sum       []byte
	CreatedAt time.Time `gorm:"<-:create"`
}

// TableName ...
func (d *ScriptVersion) TableName() string {
	return "script_versions"
}

// Delete ...
func (n ScriptVersions) Delete(ctx context.Context, id int64) (err error) {
	if err = n.Db.WithContext(ctx).Delete(&ScriptVersion{Id: id}).Error; err != nil {
		err = errors.Wrap(apperr.ErrScriptVersionDelete, err.Error())
	}
	return
}
