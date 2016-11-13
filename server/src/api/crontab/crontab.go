package crontab

import "log"

// Singleton
var instantiated *Crontab = nil

func CrontabPtr() *Crontab {
	return instantiated
}

type Crontab struct {
	Tasks []*Task
}

func (c *Crontab) NewTask(t string, h func()) *Task {
	c.Tasks = append(c.Tasks, &Task{_time:t, _func:h})
	return c.Tasks[len(c.Tasks) - 1].Run()
}

func (c *Crontab) Run() {
	for _, task := range c.Tasks {
		task.Run()
	}
}

func (c *Crontab) Stop() {
	for _, task := range c.Tasks {
		task.Stop()
	}
}

func Initialize() (err error) {
	log.Println("Crontab initialize...")

	if instantiated == nil {
		instantiated = &Crontab{
			Tasks: make([]*Task, 0),
		}
	}

	return
}