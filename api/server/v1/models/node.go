package models

import "time"

type NewNode struct {
	Port        int64  `json:"port"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type UpdateNode struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Port        int64  `json:"port"`
	Status      string `json:"status"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type Node struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        int       `json:"port"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Nodes []*Node

type NodeListModel struct {
	Items []Node `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}
