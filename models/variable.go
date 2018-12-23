package models

import "time"

type Variable struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Autoload  bool      `json:"autoload"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
