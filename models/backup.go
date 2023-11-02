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

package models

import (
	"io/fs"
	"time"
)

type Backup struct {
	ModTime  time.Time
	Name     string
	MimeType string
	Size     int64
	FileMode fs.FileMode
}

type Backups []*Backup

func (l Backups) Len() int           { return len(l) }
func (l Backups) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Backups) Less(i, j int) bool { return l[i].ModTime.UnixNano() > l[j].ModTime.UnixNano() }
