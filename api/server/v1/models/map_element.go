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
	Name          string                  `json:"name" valid:"Required"`
	Description   string                  `json:"description"`
	PrototypeId   int64                   `json:"prototype_id"`
	PrototypeType string                  `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype" valid:"Required"`
	MapId         int64                   `json:"map_id" valid:"Required"`
	LayerId       int64                   `json:"layer_id" valid:"Required"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        string                  `json:"status" valid:"Required"`
	Weight        int                     `json:"weight"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}
