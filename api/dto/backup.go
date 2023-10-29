// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package dto

import (
	"github.com/e154/smart-home/api/stub"
	m "github.com/e154/smart-home/models"
)

// Backup ...
type Backup struct{}

// NewBackupDto ...
func NewBackupDto() Backup {
	return Backup{}
}

func (b *Backup) ToBackupListResult(images []*m.Backup) []*stub.ApiBackup {

	var items = make([]*stub.ApiBackup, 0, len(images))
	for _, item := range images {
		items = append(items, &stub.ApiBackup{
			FileMode: uint32(item.FileMode),
			ModTime:  item.ModTime,
			Name:     item.Name,
			Size:     item.Size,
		})
	}

	return items
}
