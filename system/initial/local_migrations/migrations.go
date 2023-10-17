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

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/logger"
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
		var ok []string

		defer func() {
			ch <- err
			close(ch)
			if len(ok) > 0 {
				fmt.Println("\n\r")
				for _, item := range ok {
					log.Infof("migration '%s' ... installed", item)
				}
			}
			log.Infof("Applied %d migrations!", len(ok))
		}()

		var list []string
		var position = 0
		if ver != "" {
			for i, migration := range t.list {
				name := reflect.TypeOf(migration).String()
				list = append(list, name)
				if ver == name {
					position = i
				}
			}
		}

		if position >= len(t.list)-1 {
			return
		}

		for _, migration := range t.list[position+1 : len(t.list)] {
			if err = ctx.Err(); err != nil {
				log.Error(err.Error())
				return
			}

			newVersion = reflect.TypeOf(migration).String()
			tx := adaptors.Begin()
			if err = migration.Up(ctx, tx); err != nil {
				fmt.Printf("\n\nmigration '%s' ... error\n", newVersion)
				log.Error(err.Error())
				_ = tx.Rollback()
				return
			}
			_ = tx.Commit()
			ok = append(ok, newVersion)
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
