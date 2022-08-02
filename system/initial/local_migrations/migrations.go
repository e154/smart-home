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
	ctxTimeout = 5
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

func (t *Migrations) Up(ctx context.Context, tx *adaptors.Adaptors, ver string) (newVersion string, err error) {

	ctx, ctxCancel := context.WithTimeout(ctx, time.Second*ctxTimeout)
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

		for _, migration := range t.list[position:len(t.list)] {
			if err = ctx.Err(); err != nil {
				return
			}
			newVersion = reflect.TypeOf(migration).String()
			if err = migration.Up(ctx, tx); err != nil {
				fmt.Printf("migration '%s' ... error\n", newVersion)
				return
			}
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
