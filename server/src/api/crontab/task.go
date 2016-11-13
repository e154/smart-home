package crontab

import (
	"time"
)

type Task struct {
	_time	string
	_func	func()
	quit	chan bool
	isRun	bool
}

//TODO need normal cron like runner
func (t *Task) Run() *Task {

	if t.isRun {
		return t
	}

	ticker := time.NewTicker(10 * time.Second)
	t.quit = make(chan bool)
	t.isRun = true
	go func() {
		for {
			select {
			case <- ticker.C:
			// do stuff
				t._func()

			case <- t.quit:
				ticker.Stop()
				close(t.quit)
				t.isRun = false
				return
			}
		}
	}()

	return t
}

func (t *Task) Stop() *Task {
	if !t.isRun {
		return t
	}

	t.quit <- true
	return t
}

func (t *Task) IsRun() bool {
	return t.isRun
}

func (t *Task) Restart() {
	t.Stop()
	t.Run()
}