package debounce

import (
	"sync"
	"time"
)

func New(after time.Duration) func(f func()) {
	d := &Debounce{after: after}

	return func(f func()) {
		d.add(f)
	}
}

type Debounce struct {
	sync.Mutex
	after time.Duration
	timer *time.Timer
}

func (d *Debounce) add(f func()) {
	d.Lock()
	defer d.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.after, f)
}
