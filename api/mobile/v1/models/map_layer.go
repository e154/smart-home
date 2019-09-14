package models

// swagger:model
type MapLayer struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Map         *Map          `json:"map"`
	MapId       int64         `json:"map_id"`
	Status      string        `json:"status"`
	Weight      int64         `json:"weight"`
	Elements    []*MapElement `json:"elements"`
}

// swagger:model
type NewMapLayer struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Map         *Map   `json:"map"`
	Status      string `json:"status"`
}

// swagger:model
type UpdateMapLayer struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Map         *Map   `json:"map"`
	Status      string `json:"status"`
}

type SortMapLayer struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}
