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

package plugins

import (
	"fmt"
	. "github.com/e154/smart-home/tests/plugins/container"
	"go.uber.org/dig"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func init() {
	apppath := filepath.Join(os.Getenv("PWD"), "../..")
	os.Chdir(apppath)
}

var (
	container *dig.Container
)

func TestMain(m *testing.M) {

	runtime.GOMAXPROCS(-1)

	container = BuildContainer()
	err := container.Invoke(func(
		//logging *logging.Logging,
		) {

		time.Sleep(time.Millisecond * 500)

		os.Exit(m.Run())
	})

	if err != nil {
		fmt.Println("error:", dig.RootCause(err))
	}
}
