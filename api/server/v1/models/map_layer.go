package models

import "time"

type MapLayer struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Map         *Map          `json:"map"`
	MapId       int64         `json:"map_id"`
	Status      string        `json:"status"`
	Weight      int64         `json:"weight"`
	Elements    []*MapElement `json:"elements"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
