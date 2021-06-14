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

package cron

import (
	"strings"
	"sync"
)

// Task ...
type Task struct {
	_time map[int][]int
	_func func()
	cron  *Cron
	sync.Mutex
	enabled bool
}

// Enable ...
func (t *Task) Enable() *Task {
	t.Lock()
	t.enabled = true
	t.Unlock()
	return t
}

// Disable ...
func (t *Task) Disable() *Task {
	t.Lock()
	t.enabled = false
	t.Unlock()
	return t
}

// Enabled ...
func (t *Task) Enabled() bool {
	t.Lock()
	defer t.Unlock()
	return t.enabled
}

// SetTime ...
func (t *Task) SetTime(time string) {
	args := strings.Split(time, " ")
	switch len(args) {
	case 1:
		if args[0] == "*" {

		}

	}
}

// GetTime ...
func (t *Task) GetTime() (time string) {

	return
}

func (t *Task) exec(_timer *Timer) {

	//log.Println("exec")

	// WEEKDAY
	exist := false
	for _, weekday := range t._time[WEEKDAY] {
		if weekday == int(_timer.weekday) {
			exist = true
			break
		}
	}

	if !exist {
		return
	}

	//log.Println("weekday")

	// MONTH
	exist = false
	for _, month := range t._time[MONTH] {
		if month == int(_timer.month) {
			exist = true
			break
		}
	}
	if !exist {
		return
	}

	//log.Println("month")

	// DAY
	exist = false
	for _, day := range t._time[DAY] {
		if day == _timer.day {
			exist = true
			break
		}
	}
	if !exist {
		return
	}

	//log.Println("day")

	// HOUR
	exist = false
	for _, hour := range t._time[HOUR] {
		if hour == _timer.hour {
			exist = true
			break
		}
	}
	if !exist {
		return
	}

	//log.Println("hour")

	// MINUTES
	exist = false
	for _, min := range t._time[MINUTE] {
		if min == _timer.min {
			exist = true
			break
		}
	}
	if !exist {
		return
	}

	//log.Println("minutes")

	// SECONDS
	exist = false
	for _, sec := range t._time[SECOND] {
		if sec == _timer.second {
			exist = true
			break
		}
	}
	if !exist {
		return
	}

	//log.Println("seconds")

	t.Run()
}

// Run ...
func (t *Task) Run() {
	t._func()
}
