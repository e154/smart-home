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

package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
)

var Restore bool

// Start ...
func Start(app *fx.App) {
	startCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}
}

// Work ...
func Work() {
	var gracefulStop = make(chan os.Signal, 10)
	signal.Notify(gracefulStop, syscall.SIGINT, syscall.SIGTERM)

	<-gracefulStop
}

// Stop ...
func Stop(app *fx.App) {
	t := 60 * time.Second
	if Restore {
		t = 15 * time.Minute
	}
	stopCtx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Kill() error {
	return syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func Do(builder func(opt fx.Option) (app *fx.App), options fx.Option) {
	app := builder(options)

	Start(app)

	Work()

	Stop(app)
}
