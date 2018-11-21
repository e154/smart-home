package cron

import "time"

type Timer struct {
	second    int
	min       int
	hour      int
	weekday   time.Weekday
	day       int
	month     time.Month
}