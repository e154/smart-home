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

package api

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/e154/smart-home/internal/system/validation"

	"go.uber.org/dig"

	. "github.com/e154/smart-home/tests/api/container"
)

func init() {
	apppath := filepath.Join(os.Getenv("PWD"), "../..")
	_ = os.Chdir(apppath)
}

var (
	container *dig.Container
)

func TestMain(m *testing.M) {

	runtime.GOMAXPROCS(-1)

	_ = os.Setenv("TEST_MODE", "true")

	container = BuildContainer()
	err := container.Invoke(func(
		validation *validation.Validate,
		//logging *logging.Logging,
	) {
		_ = validation.Start(context.Background())
		time.Sleep(time.Millisecond * 500)
		os.Exit(m.Run())
	})

	if err != nil {
		fmt.Println("error:", dig.RootCause(err))
	}
}
