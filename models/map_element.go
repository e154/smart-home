package models

import (
	"encoding/json"
	"time"
	"github.com/e154/smart-home/system/validation"
	. "github.com/e154/smart-home/common"
)

type MapElement struct {
	Id            int64           `json:"id"`
	Name          string          `json:"name" valid:"Required"`
	Description   string          `json:"description"`
	PrototypeId   int64           `json:"prototype_id"`
	PrototypeType PrototypeType   `json:"prototype_type"`
	Prototype     interface{}     `json:"prototype" valid:"Required"`
	MapId         int64           `json:"map_id" valid:"Required"`
	LayerId       int64           `json:"layer_id" valid:"Required"`
	GraphSettings json.RawMessage `json:"graph_settings"`
	Status        StatusType      `json:"status" valid:"Required"`
	Weight        int             `json:"weight"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (m *MapElement) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}
