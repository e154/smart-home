package models

import (
	"time"
)

type MapOptions struct {
	Zoom              float64 `json:"zoom"`
	ElementStateText  bool    `json:"element_state_text"`
	ElementOptionText bool    `json:"element_option_text"`
}

type NewMap struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

type UpdateMap struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
}

type Map struct {
	Id          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Options     MapOptions `json:"options"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type MapFullModel struct {
	Id          int64       `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Options     MapOptions  `json:"options"`
	Layers      []*MapLayer `json:"layers"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Maps []*Map

type MapListModel struct {
	Items []Map `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}

type SearchMapResponse struct {
	Maps []Map `json:"maps"`
}
