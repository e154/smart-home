package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
	"errors"
)

func AddMap(m *m.Map, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation
	ok, errs = m.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if id, err = adaptors.Map.Add(m); err != nil {
		return
	}

	m.Id = id


	return
}

func GetMapById(mId int64, adaptors *adaptors.Adaptors) (m *m.Map, err error) {

	m, err = adaptors.Map.GetById(mId)

	return
}

func UpdateMap(mapParams *models.UpdateMap, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	var m *m.Map
	if m, err = adaptors.Map.GetById(mapParams.Id); err != nil {
		return
	}

	copier.Copy(&m, &mapParams)

	// validation
	ok, errs = m.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	if err = adaptors.Map.Update(m); err != nil {
		return
	}

	return
}

func GetMapList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Map, total int64, err error) {

	items, total, err = adaptors.Map.List(limit, offset, order, sortBy)

	return
}

func DeleteMapById(mId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	var m *m.Map
	if m, err = adaptors.Map.GetById(mId); err != nil {
		return
	}

	err = adaptors.Map.Delete(m.Id)

	return
}

func SearchMap(query string, limit, offset int, adaptors *adaptors.Adaptors) (nodes []*m.Map, total int64, err error) {

	nodes, total, err = adaptors.Map.Search(query, limit, offset)

	return
}
