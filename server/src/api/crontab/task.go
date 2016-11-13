package crontab

import "time"

type Task struct {
	_time	string
	_func	func()
	quit	chan struct{}
	isRun	bool
}

//TODO need normal cron like runner
func (t *Task) Run() *Task {

	ticker := time.NewTicker(10 * time.Second)
	t.quit = make(chan struct{})
	go func() {
		t.isRun = true
		for {
			select {
			case <- ticker.C:
			// do stuff
				t._func()

			case <- t.quit:
				ticker.Stop()
				t.isRun = false
				return
			}
		}
	}()

	return t
}

func (t *Task) Stop() *Task {
	close(t.quit)
	return t
}

func (t *Task) IsRun() bool {
	return t.isRun
}

func (t *Task) Restart() {
	t.Stop()
	t.Run()
}