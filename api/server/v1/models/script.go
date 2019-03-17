package models

import (
	"time"
)

// swagger:model
type NewScript struct {
	Lang        string `json:"lang"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// swagger:model
type UpdateScript struct {
	Id          int64  `json:"id"`
	Lang        string `json:"lang"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// swagger:model
type ExecScript struct {
	Lang        string `json:"lang"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

// swagger:model
type Script struct {
	Id          int64     `json:"id"`
	Lang        string    `json:"lang"`
	Name        string    `json:"name"`
	Source      string    `json:"source"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
