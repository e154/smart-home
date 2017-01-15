package models

import "time"

type UserHistory struct {
	Ip	string		`json:"ip"`
	Time	time.Time	`json:"time"`
}
