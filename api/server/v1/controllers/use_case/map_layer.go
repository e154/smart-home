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

func AddMapLayer(params *models.NewMapLayer,
	adaptors *adaptors.Adaptors) (result *models.MapLayer, id int64, errs []*validation.Error, err error) {

	m := &m.MapLayer{}
	common.Copy(&m, &params)

	if params.Map != nil && params.Map.Id != 0 {
		m.MapId = params.Map.Id
	}

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = adaptors.MapLayer.Add(m); err != nil {
		return
	}

	if m, err = adaptors.MapLayer.GetById(id); err != nil {
		return
	}

	result = &models.MapLayer{}
	err = common.Copy(&result, &m)

	return
}

func GetMapLayerById(mId int64, adaptors *adaptors.Adaptors) (result *models.MapLayer, err error) {

	var m *m.MapLayer
	if m, err = adaptors.MapLayer.GetById(mId); err != nil {
		return
	}

	result = &models.MapLayer{}
	err = common.Copy(&result, &m)

	return
}

func UpdateMapLayer(mapParams *models.UpdateMapLayer, adaptors *adaptors.Adaptors) (ok bool, errs []*validation.Error, err error) {

	var m *m.MapLayer
	if m, err = adaptors.MapLayer.GetById(mapParams.Id); err != nil {
		return
	}

	copier.Copy(&m, &mapParams)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	err = adaptors.MapLayer.Update(m)

	return
}

func SortMapLayers(params []*models.SortMapLayer, adaptors *adaptors.Adaptors) (err error) {

	for _, s := range params {
		adaptors.MapLayer.Sort(&m.MapLayer{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

func DeleteMapLayer(mId int64, adaptors *adaptors.Adaptors) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	if _, err = adaptors.MapLayer.GetById(mId); err != nil {
		return
	}

	err = adaptors.MapLayer.Delete(mId)

	return
}

func GetMapLayerList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (result []*models.MapLayer, total int64, err error) {

	var items []*m.MapLayer
	if items, total, err = adaptors.MapLayer.List(limit, offset, order, sortBy); err != nil {
		return
	}

	result = make([]*models.MapLayer, 0)
	err = common.Copy(&result, &items)

	return
}
