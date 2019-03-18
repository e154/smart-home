package use_case

import (
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/validation"
	m "github.com/e154/smart-home/models"
	"github.com/jinzhu/copier"
	"errors"
)

func AddMapElement(params *models.NewMapElement,
	adaptors *adaptors.Adaptors) (result *models.MapElement, id int64, errs []*validation.Error, err error) {

	mapElement := &m.MapElement{}
	common.Copy(&mapElement, &params, common.JsonEngine)

	if params.Map.Id != 0 {
		mapElement.MapId = params.Map.Id
	}

	if params.Layer.Id != 0 {
		mapElement.LayerId = params.Layer.Id
	}

	// validation
	_, errs = mapElement.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = adaptors.MapElement.Add(mapElement); err != nil {
		return
	}

	if mapElement, err = adaptors.MapElement.GetById(id); err != nil {
		return
	}

	result = &models.MapElement{}
	err = common.Copy(&result, &mapElement, common.JsonEngine)

	return
}

func GetMapElementById(mId int64, adaptors *adaptors.Adaptors) (result *models.MapElement, err error) {

	var m *m.MapElement
	if m, err = adaptors.MapElement.GetById(mId); err != nil {
		return
	}

	result = &models.MapElement{}
	err = common.Copy(&result, &m)

	return
}

func UpdateMapElement(mapParams *models.UpdateMapElement, adaptors *adaptors.Adaptors) (result *models.MapElement, errs []*validation.Error, err error) {

	var m *m.MapElement
	if m, err = adaptors.MapElement.GetById(mapParams.Id); err != nil {
		return
	}

	copier.Copy(&m, &mapParams)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.MapElement.Update(m); err != nil {
		return
	}

	if m, err = adaptors.MapElement.GetById(m.Id); err != nil {
		return
	}

	result = &models.MapElement{}
	err = common.Copy(&result, &m)

	return
}

func SortMapElements(params []*models.SortMapElement, adaptors *adaptors.Adaptors) (err error) {

	for _, s := range params {
		adaptors.MapElement.Sort(&m.MapElement{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

func DeleteMapElement(mId int64, adaptors *adaptors.Adaptors) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	if _, err = adaptors.MapElement.GetById(mId); err != nil {
		return
	}

	err = adaptors.MapElement.Delete(mId)

	return
}

func GetMapElementList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.MapElement, total int64, err error) {

	var items []*m.MapElement
	if items, total, err = adaptors.MapElement.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.MapElement, 0)
	err = common.Copy(&result, &items)

	return
}
