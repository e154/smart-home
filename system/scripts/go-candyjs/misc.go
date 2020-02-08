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

package candyjs

import "strings"

func isExported(name string) bool {
	return nameToJavaScript(name) != name
}

func nameToJavaScript(name string) string {
	var toLower, keep string
	for _, c := range name {
		if c >= 'A' && c <= 'Z' && len(keep) == 0 {
			toLower += string(c)
		} else {
			keep += string(c)
		}
	}

	lc := len(toLower)
	if lc > 1 && lc != len(name) {
		keep = toLower[lc-1:] + keep
		toLower = toLower[:lc-1]

	}

	return strings.ToLower(toLower) + keep
}

func nameToGo(name string) []string {
	if name[0] >= 'A' && name[0] <= 'Z' {
		return nil
	}

	var toUpper, keep string
	for _, c := range name {
		if c >= 'a' && c <= 'z' && len(keep) == 0 {
			toUpper += string(c)
		} else {
			keep += string(c)
		}
	}

	return []string{
		strings.Title(name),
		strings.ToUpper(toUpper) + keep,
	}
}
