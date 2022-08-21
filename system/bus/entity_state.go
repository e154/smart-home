package bus

import (
	"time"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
)

// EntityState ...
type EntityState struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageUrl    *string `json:"image_url"`
	Icon        *string `json:"icon"`
}

// EventEntityState ...
type EventEntityState struct {
	EntityId    common.EntityId `json:"entity_id"`
	Value       interface{}     `json:"value"`
	State       *EntityState    `json:"state"`
	Attributes  m.Attributes    `json:"attributes"`
	Settings    m.Attributes    `json:"settings"`
	LastChanged *time.Time      `json:"last_changed"`
	LastUpdated *time.Time      `json:"last_updated"`
}

// Compare ...
func (e1 EventEntityState) Compare(e2 EventEntityState) (ident bool) {

	if e1.State == nil || e2.State == nil {
		return
	}

	if e1.State.Name != e2.State.Name {
		return
	}

	for k1, v1 := range e1.Attributes {
		if !e2.Attributes[k1].Compare(v1) {
			return false
		}
	}

	ident = true

	return
}
