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

package demo

import (
	"context"
	"time"

	"github.com/e154/smart-home/pkg/logger"
)

const (
	ctxTimeout = 5
)

var (
	log = logger.MustGetLogger("demo")
)

type Demos struct {
	list map[string]Demo
}

func NewDemos(list map[string]Demo) *Demos {
	return &Demos{
		list: list,
	}
}

func (t *Demos) InstallByName(ctx context.Context, name string) (err error) {

	if name == "" {
		return
	}

	ctx, ctxCancel := context.WithTimeout(ctx, time.Second*ctxTimeout)
	defer ctxCancel()

	log.Infof("install demo \"%s\" ...", name)

	ch := make(chan error, 1)

	go func() {
		var err error
		defer func() {
			ch <- err
			close(ch)
		}()

		if err = ctx.Err(); err != nil {
			return
		}
		if err = t.list[name].Install(ctx); err != nil {
			return
		}
	}()

	select {
	case v := <-ch:
		err = v
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}
