package demo

import (
	"context"
	"time"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/common/logger"
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

func (t *Demos) InstallByName(ctx context.Context, adaptors *adaptors.Adaptors, name string) (err error) {

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
		if err = t.list[name].Install(ctx, adaptors); err != nil {
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
