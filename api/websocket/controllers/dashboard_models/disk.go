// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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

//+build linux windows darwin,!386

package dashboard_models

import (
	"github.com/shirou/gopsutil/disk"
	"sync"
)

func NewDisk() (_disk *Disk) {

	var state_root *disk.UsageStat
	//var state_tmp *disk.UsageStat
	var err error

	_disk = &Disk{}

	if state_root, err = disk.Usage("/"); err == nil {
		_disk.Root = state_root
	}

	//if state_tmp, err = disk.Usage("/tmp"); err == nil {
	//	_disk.Tmp = state_tmp
	//}

	return
}

type Disk struct {
	sync.Mutex
	Root *disk.UsageStat `json:"root"`
	//Tmp		*disk.UsageStat		`json:"tmp"`
}

func (d *Disk) Update() {

	var state_root *disk.UsageStat
	//var state_tmp *disk.UsageStat
	var err error

	if state_root, err = disk.Usage("/"); err == nil {
		d.Lock()
		d.Root = state_root
		d.Unlock()
	}

	//if state_tmp, err = disk.Usage("/tmp"); err == nil {
	//	d.Tmp = state_tmp
	//}
}
