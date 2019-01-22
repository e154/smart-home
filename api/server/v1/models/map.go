package models

import (
	"encoding/json"
	"time"
)

type NewMap struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     json.RawMessage `json:"options"`
}

type UpdateMap struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     json.RawMessage `json:"options"`
}

type Map struct {
	Id          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Options     json.RawMessage `json:"options"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
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
	Maps []Map `json:"nodes"`
}
