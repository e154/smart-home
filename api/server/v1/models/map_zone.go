package models

type NewMapZone struct {
	Name string `json:"name"`
}

// swagger:model
type MapZone struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
