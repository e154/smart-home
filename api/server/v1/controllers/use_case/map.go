package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
	"errors"
	"github.com/e154/smart-home/common"
)

func AddMap(params *models.NewMap, adaptors *adaptors.Adaptors, core *core.Core) (result *models.Map, id int64, errs []*validation.Error, err error) {

	m := &m.Map{}
	common.Copy(&m, &params)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = adaptors.Map.Add(m); err != nil {
		return
	}

	if m, err = adaptors.Map.GetById(id); err != nil {
		return
	}

	result = &models.Map{}
	err = common.Copy(&result, &m)

	return
}

func GetMapById(mId int64, adaptors *adaptors.Adaptors) (result *models.Map, err error) {

	var m *m.Map
	if m, err = adaptors.Map.GetById(mId); err != nil {
		return
	}

	result = &models.Map{}
	err = common.Copy(&result, &m)

	return
}

func GetFullMapById(mId int64, adaptors *adaptors.Adaptors) (result *models.MapFullModel, err error) {

	var m *m.Map
	if m, err = adaptors.Map.GetFullById(mId); err != nil {
		return
	}

	//debug.Println(m)
	//fmt.Println("------------------------")

	result = &models.MapFullModel{}
	err = common.Copy(&result, &m, common.JsonEngine)

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

	err = adaptors.Map.Update(m)

	return
}

func GetMapList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.Map, total int64, err error) {

	var items []*m.Map
	if items, total, err = adaptors.Map.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.Map, 0)
	err = common.Copy(&result, &items)

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
