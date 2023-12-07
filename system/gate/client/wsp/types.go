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

package wsp

import "fmt"

// PoolSize represent the number of open connections per status
type PoolSize struct {
	connecting int
	idle       int
	running    int
	total      int
}

type Status struct {
	Connecting int
	Idle       int
	Running    int
	Total      int
}

func (poolSize *PoolSize) String() string {
	return fmt.Sprintf("Connecting %d, idle %d, running %d, total %d", poolSize.connecting, poolSize.idle, poolSize.running, poolSize.total)
}
