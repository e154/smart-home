package models

import "time"

type NewNode struct {
	Port        int64  `json:"port"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type UpdateNodeModel struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Port        int64  `json:"port"`
	Status      string `json:"status"`
	IP          string `json:"ip"`
	Description string `json:"description"`
}

type NodeModel struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        int       `json:"port"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Nodes []*NodeModel

type NodeListModel struct {
	Items []NodeModel `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}

type ResponseNodeModel struct {
	Node NodeModel `json:"node"`
}

type ResponseSearchNode struct {
	Nodes []NodeModel `json:"nodes"`
}
