package models

import "time"

// swagger:model
type NewLog struct {
	Body  string `json:"body"`
	Level string `json:""`
}

// swagger:model
type Log struct {
	Id        int64     `json:"id"`
	Body      string    `json:"body"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}
