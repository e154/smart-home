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

package debug

import (
	"encoding/json"
	"fmt"
	"github.com/e154/smart-home/system/validation"
	"runtime"
	"strings"
)

// https://github.com/Unknwon/gcblog/blob/master/content/04-go-caller.md
func CallerName(skip int) (name, file string, line int, ok bool) {
	var pc uintptr
	if pc, file, line, ok = runtime.Caller(skip + 1); !ok {
		return
	}
	name = runtime.FuncForPC(pc).Name()
	return
}

// Trace ...
func Trace() (trace string) {

	i := 1 //0...
	for skip := i; ; skip++ {
		name, file, line, ok := CallerName(skip)
		if !ok {
			break
		}
		fn := strings.Title(strings.Split(name, ".")[1]) + "()"
		trace += "\n"
		trace += fmt.Sprintf("called: %s:%s line: %d", file, fn, line)
	}

	return
}

// Println ...
func Println(i interface{}) {
	b, _ := json.MarshalIndent(i, " ", "  ")
	fmt.Println(string(b))
}

// PrintValidationErrs ...
func PrintValidationErrs(errs []*validation.Error) {
	for _, err := range errs {
		fmt.Printf("%s - %s", err.Name, err.String())
	}
}
