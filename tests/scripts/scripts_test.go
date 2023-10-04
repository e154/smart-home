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

package scripts

import (
	"encoding/hex"
	"fmt"
	"github.com/e154/smart-home/common/encryptor"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/e154/smart-home/tests/scripts/container"
	"go.uber.org/dig"
)

func init() {
	apppath := filepath.Join(os.Getenv("PWD"), "../..")
	_ = os.Chdir(apppath)
}

var (
	container *dig.Container
)

func TestMain(m *testing.M) {

	container = BuildContainer()
	err := container.Invoke(func() {

		// encryptor
		b, _ := hex.DecodeString("7abf835e883087d3dc87be2c24ea2faee948f03cf28ebe6d7c119c2ccedc9ab2")
		encryptor.SetKey(b)

		time.Sleep(time.Millisecond * 500)

		os.Exit(m.Run())
	})

	if err != nil {
		fmt.Println("error:", dig.RootCause(err))
	}
}
