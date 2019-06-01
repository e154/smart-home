package endpoint

import (
	"errors"
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common"
)

type MapCommand struct {
	*CommonCommand
}

func NewMapCommand(common *CommonCommand) *MapCommand {
	return &MapCommand{
		CommonCommand: common,
	}
}

func (m *MapCommand) Add(params *m.Map) (result *m.Map, errs []*validation.Error, err error) {

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

func (m *MapCommand) GetById(id int64) (result *m.Map, err error) {

	result, err = m.adaptors.Map.GetById(id)

	return
}

func (m *MapCommand) GetFullById(mId int64) (result *m.Map, err error) {

	result, err = m.adaptors.Map.GetFullById(mId)

	return
}

func (n *MapCommand) Update(params *m.Map) (result *m.Map, errs []*validation.Error, err error) {

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

func (n *MapCommand) GetList(limit, offset int64, order, sortBy string) (items []*m.Map, total int64, err error) {

	items, total, err = n.adaptors.Map.List(limit, offset, order, sortBy)

	return
}

func (n *MapCommand) Delete(mId int64) (err error) {

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

func (n *MapCommand) Search(query string, limit, offset int) (items []*m.Map, total int64, err error) {

	items, total, err = n.adaptors.Map.Search(query, limit, offset)

	return
}
