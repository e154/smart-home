package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/null"
)

type MapElementGraphSettingsPosition struct {
	Top  int64 `json:"top"`
	Left int64 `json:"left"`
}
type MapElementGraphSettings struct {
	Width    null.Int64                      `json:"width"`
	Height   null.Int64                      `json:"height"`
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
	PrototypeType PrototypeType           `json:"prototype_type"`
	Prototype     Prototype               `json:"prototype" valid:"Required"`
	MapId         int64                   `json:"map_id" valid:"Required"`
	LayerId       int64                   `json:"layer_id" valid:"Required"`
	GraphSettings MapElementGraphSettings `json:"graph_settings"`
	Status        StatusType              `json:"status" valid:"Required"`
	Weight        int64                   `json:"weight"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
}

func (m *MapElement) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}

type SortMapElementByWeight []*MapElement

func (l SortMapElementByWeight) Len() int           { return len(l) }
func (l SortMapElementByWeight) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SortMapElementByWeight) Less(i, j int) bool { return l[i].Weight < l[j].Weight }
