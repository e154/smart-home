package models

import (
	"time"
)

type MapElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}
type MapElementGraphSettings struct {
	Width    *int64                          `json:"width"`
	Height   *int64                          `json:"height"`
	Position MapElementGraphSettingsPosition `json:"position"`
}

type Prototype struct {
	*MapImage
	*MapText
	*MapDevice
}

type MapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

type NewMapElement struct {
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	Map           *Map                    `json:"map"`
	Layer         *MapLayer               `json:"layer"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
}

type UpdateMapElement struct {
	Id            int64                   `json:"id"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype"`
	Map           *Map                    `json:"map"`
	MapId         int64                   `json:"map_id"`
	LayerId       int64                   `json:"layer_id"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status"`
	Weight        int                     `json:"weight"`
}

type SortMapElement struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}
