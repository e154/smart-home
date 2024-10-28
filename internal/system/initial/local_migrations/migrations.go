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

package local_migrations

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/logger"
)

const (
	ctxTimeout = 30 * time.Second
)

var (
	log = logger.MustGetLogger("local_migrations")
)

type Migrations struct {
	list []Migration
}

func NewMigrations(list []Migration) *Migrations {
	return &Migrations{
		list: list,
	}
}

func (t *Migrations) Up(ctx context.Context, adaptors *adaptors.Adaptors, ver string) (newVersion string, err error) {

	ctx, ctxCancel := context.WithTimeout(ctx, ctxTimeout)
	defer ctxCancel()

	ch := make(chan error, 1)

	newVersion = ver

	go func() {
		var err error
		var done []string

		defer func() {
			ch <- err
			close(ch)
			if len(done) > 0 {
				fmt.Println("\n\r")
				for _, item := range done {
					log.Infof("migration '%s' ... installed", item)
				}
			}
			log.Infof("Applied %d migrations!", len(done))
		}()

		var list []string
		var position = 0
		var exist = false
		if ver != "" {
			for i, migration := range t.list {
				name := reflect.TypeOf(migration).String()
				list = append(list, name)
				if ver == name {
					exist = true
					position = i
					break
				}
			}
		}

		if !exist {
			log.Errorf("Unknown migration %s!", ver)
			err = fmt.Errorf(fmt.Sprintf("Unknown migration %s!", ver))
			return
		}

		if position >= len(t.list)-1 {
			return
		}

		if position > 0 {
			position++
		}

		for _, migration := range t.list[position:len(t.list)] {
			if err = ctx.Err(); err != nil {
				log.Error(err.Error())
				return
			}

			itemVersion := reflect.TypeOf(migration).String()

			err = adaptors.Transaction.Do(ctx, func(ctx context.Context) error {
				return migration.Up(ctx)
			})
			if err != nil {
				log.Errorf("migration '%s' ended with error: %s", itemVersion, err.Error())
				break
			}
			newVersion = itemVersion
			done = append(done, itemVersion)
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
