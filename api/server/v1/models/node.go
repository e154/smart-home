package models

import "time"

// swagger:model
type NewNode struct {
	Port           int64  `json:"port"`
	Status         string `json:"status"`
	Name           string `json:"name"`
	IP             string `json:"ip"`
	Description    string `json:"description"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

// swagger:model
type UpdateNode struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Port           int64  `json:"port"`
	Status         string `json:"status"`
	IP             string `json:"ip"`
	Description    string `json:"description"`
	Login          string `json:"login"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

// swagger:model
type Node struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Ip          string    `json:"ip"`
	Port        int       `json:"port"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Login       string    `json:"login"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
