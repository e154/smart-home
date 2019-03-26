package models

import (
	"time"
	"github.com/e154/smart-home/system/validation"
)

type MapLayer struct {
	Id          int64         `json:"id"`
	Name        string        `json:"name" valid:"MaxSize(254);Required"`
	Description string        `json:"description"`
	Map         *Map          `json:"map"`
	MapId       int64         `json:"map_id" valid:"Required"`
	Status      string        `json:"status" valid:"Required"`
	Weight      int64         `json:"weight"`
	Elements    []*MapElement `json:"elements"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (m *MapLayer) Valid() (ok bool, errs []*validation.Error) {

	valid := validation.Validation{}
	if ok, _ = valid.Valid(m); !ok {
		errs = valid.Errors
	}

	return
}

type SortMapLayersByWeight []*MapLayer

func (l SortMapLayersByWeight) Len() int           { return len(l) }
func (l SortMapLayersByWeight) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l SortMapLayersByWeight) Less(i, j int) bool { return l[i].Weight < l[j].Weight }

type SortMapLayer struct {
	Id     int64 `json:"id"`
	Weight int64 `json:"weight"`
}
