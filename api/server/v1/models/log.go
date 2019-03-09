package models

import "time"

type NewLogModel struct {
	Body  string `json:"body"`
	Level string `json:""`
}

type LogModel struct {
	Id        int64     `json:"id"`
	Body      string    `json:"body"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

type ResponseLog struct {
	Log *LogModel `json:"log"`
}

type ResponseSearchLog struct {
	Logs []*LogModel `json:"logs"`
}

type ResponseLogList struct {
	Items []LogModel `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}
