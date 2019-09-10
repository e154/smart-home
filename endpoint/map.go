package endpoint

import (
	"errors"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

type MapEndpoint struct {
	*CommonEndpoint
}

func NewMapEndpoint(common *CommonEndpoint) *MapEndpoint {
	return &MapEndpoint{
		CommonEndpoint: common,
	}
}

func (m *MapEndpoint) Add(params *m.Map) (result *m.Map, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = m.adaptors.Map.Add(params); err != nil {
		return
	}

	result, err = m.adaptors.Map.GetById(id)

	return
}

func (m *MapEndpoint) GetById(id int64) (result *m.Map, err error) {

	result, err = m.adaptors.Map.GetById(id)

	return
}

func (m *MapEndpoint) GetActiveElements(sortBy, order string, limit, offset int) (result []*m.MapElement, total int64, err error) {

	result, total, err = m.adaptors.MapElement.GetActiveElements(sortBy, order, limit, offset)

	return
}

func (m *MapEndpoint) GetFullById(mId int64) (result *m.Map, err error) {

	result, err = m.adaptors.Map.GetFullById(mId)

	return
}

func (n *MapEndpoint) Update(params *m.Map) (result *m.Map, errs []*validation.Error, err error) {

	var m *m.Map
	if m, err = n.adaptors.Map.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&m, &params, common.JsonEngine)

	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.Map.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.Map.GetById(params.Id)

	return
}

func (n *MapEndpoint) GetList(limit, offset int64, order, sortBy string) (items []*m.Map, total int64, err error) {

	items, total, err = n.adaptors.Map.List(limit, offset, order, sortBy)

	return
}

func (n *MapEndpoint) Delete(mId int64) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	var m *m.Map
	if m, err = n.adaptors.Map.GetById(mId); err != nil {
		return
	}

	err = n.adaptors.Map.Delete(m.Id)

	return
}

func (n *MapEndpoint) Search(query string, limit, offset int) (items []*m.Map, total int64, err error) {

	items, total, err = n.adaptors.Map.Search(query, limit, offset)

	return
}
