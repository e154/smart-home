package cron

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	SECOND int = iota
	MINUTE
	HOUR
	DAY
	MONTH
	WEEKDAY
)

func NewCron() *Cron {
	return &Cron{
		tasks: make(map[*Task]bool),
	}
}

type Cron struct {
	sync.Mutex
	tasks     map[*Task]bool
	isRun     bool
	quit      chan bool
	StartTime time.Time
	Uptime    time.Duration
}

func (c *Cron) stringParse(args string) (result []int) {

	result = make([]int, 0)
	commas := strings.Split(args, ",")
	if len(commas) != 0 {
		for _, str := range commas {
			sec := c.rangeParse(str)
			result = append(result, sec...)
		}
	} else {
		sec := c.rangeParse(args)
		result = append(result, sec...)
	}

	return
}

func (c *Cron) rangeParse(str string) (result []int) {

	// [0-9] - [0-9]
	args := strings.Split(str, "-")
	if len(args) > 1 {
		min, _ := strconv.ParseInt(args[0], 0, 10)
		max, _ := strconv.ParseInt(args[1], 0, 60)
		if len(args) != 0 {
			for i := min; i <= max; i++ {
				result = append(result, int(i))
			}
		}
	}

	// [*] / [0-9]
	//...

	// [0-9]
	if i, err := strconv.ParseInt(str, 0, 10); err == nil {
		result = append(result, int(i))
	}

	return
}

func (c *Cron) timeParser(t string) (result map[int][]int) {
	result = make(map[int][]int)
	args := strings.Split(t, " ")

	if len(args) != 6 {
		log.Println("error: bad time string")
		return
	}

	// WEEKDAY
	// ------------------------------------------------------
	if args[WEEKDAY] == "*" {
		result[WEEKDAY] = []int{}
		for i := 0; i < 7; i++ {
			result[WEEKDAY] = append(result[WEEKDAY], i)
		}

	} else {
		weekdays := c.stringParse(args[WEEKDAY])
		result[WEEKDAY] = append(result[WEEKDAY], weekdays...)
	}

	// MONTH
	// ------------------------------------------------------
	if args[MONTH] == "*" {
		result[MONTH] = []int{}
		for i := 1; i < 13; i++ {
			result[MONTH] = append(result[MONTH], i)
		}

	} else {
		months := c.stringParse(args[MONTH])
		result[MONTH] = append(result[MONTH], months...)
	}

	// DAY
	// ------------------------------------------------------
	if args[DAY] == "*" {
		result[DAY] = []int{}
		for i := 1; i < 32; i++ {
			result[DAY] = append(result[DAY], i)
		}

	} else {
		days := c.stringParse(args[DAY])
		result[DAY] = append(result[DAY], days...)
	}

	// HOUR
	// ------------------------------------------------------
	if args[HOUR] == "*" {
		result[HOUR] = []int{}
		for i := 0; i < 24; i++ {
			result[HOUR] = append(result[HOUR], i)
		}

	} else {
		hours := c.stringParse(args[HOUR])
		result[HOUR] = append(result[HOUR], hours...)
	}

	// MINUTES
	// ------------------------------------------------------
	if args[MINUTE] == "*" {
		result[MINUTE] = []int{}
		for i := 0; i < 60; i++ {
			result[MINUTE] = append(result[MINUTE], i)
		}

	} else {
		minutes := c.stringParse(args[MINUTE])
		result[MINUTE] = append(result[MINUTE], minutes...)
	}

	// SECONDS
	// ------------------------------------------------------
	if args[SECOND] == "*" {
		result[SECOND] = []int{}
		for i := 0; i < 60; i++ {
			result[SECOND] = append(result[SECOND], i)
		}

	} else {
		sec := c.stringParse(args[SECOND])
		result[SECOND] = append(result[SECOND], sec...)
	}

	return
}

func (c *Cron) NewTask(t string, h func()) *Task {
	_time := c.timeParser(t)
	task := &Task{
		_time:   _time,
		_func:   h,
		cron:    c,
		enabled: true,
	}

	c.Lock()
	c.tasks[task] = false
	c.Unlock()

	return task
}

func (c *Cron) RemoveTask(task *Task) {

	c.Lock()
	if _, ok := c.tasks[task]; !ok {
		c.Unlock()
		return
	}

	delete(c.tasks, task)
	c.Unlock()
}

func (c *Cron) timePrepare(t time.Time) {

	c.Uptime = t.Sub(c.StartTime)
	var ct *Timer = &Timer{}
	ct.second = t.Second()
	ct.min = t.Minute()
	ct.hour = t.Hour()
	ct.weekday = t.Weekday()
	ct.day = t.Day()
	ct.month = t.Month()

	// tasks
	//-----------------------------------

	c.Lock()
	defer c.Unlock()

	for task, _ := range c.tasks {
		if !task.Enabled() {
			continue
		}
		go task.exec(ct)
	}
}

func (c *Cron) Run() *Cron {
	if c.isRun {
		return c
	}

	ticker := time.NewTicker(1 * time.Second)
	c.quit = make(chan bool)
	c.isRun = true
	go func() {
		for {

			select {
			case t := <-ticker.C:
				c.timePrepare(t)
			case <-c.quit:
				ticker.Stop()
				close(c.quit)
				c.isRun = false
				return
			}
		}
	}()

	return c
}

func (c *Cron) Stop() *Cron {
	if !c.isRun {
		return c
	}

	c.quit <- true
	return c
}
