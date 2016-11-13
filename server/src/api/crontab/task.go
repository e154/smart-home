package crontab

import "time"

type Task struct {
	_time	string
	_func	func()
	quit	chan struct{}
}

//TODO need normal cron like runner
func (t *Task) Run() *Task {

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
			// do stuff
				t._func()

			case <- quit:
				ticker.Stop()
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

