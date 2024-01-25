// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package rbac

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {

	success := map[string]string{
		"/v1/backup/foo":                         "/v1/backup/[\\w]+",
		"/v1/backup/FOO":                         "/v1/backup/[\\w]+",
		"/v1/backup/123":                         "/v1/backup/[0-9]+",
		"/v1/backup/2024-01-10T11:19:49.007.zip": "/v1/backup/[\\w]+",
	}

	for s, pattern := range success {
		ok, err := regexp.MatchString(pattern, s)
		fmt.Println(pattern, s)
		require.NoError(t, err)
		require.True(t, ok)
	}
}
