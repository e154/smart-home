package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/common"
	"errors"
)

type MapLayerEndpoint struct {
	*CommonEndpoint
}

func NewMapLayerEndpoint(common *CommonEndpoint) *MapLayerEndpoint {
	return &MapLayerEndpoint{
		CommonEndpoint: common,
	}
}

func (n *MapLayerEndpoint) Add(params *m.MapLayer) (result *m.MapLayer, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.MapLayer.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.MapLayer.GetById(id)

	return
}

func (n *MapLayerEndpoint) GetById(mId int64) (result *m.MapLayer, err error) {

	result, err = n.adaptors.MapLayer.GetById(mId)

	return
}

func (n *MapLayerEndpoint) Update(params *m.MapLayer) (result *m.MapLayer, errs []*validation.Error, err error) {

	var m *m.MapLayer
	if m, err = n.adaptors.MapLayer.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&m, &params, common.JsonEngine)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	if err = n.adaptors.MapLayer.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.MapLayer.GetById(m.Id)

	return
}

func (n *MapLayerEndpoint) Sort(params []*m.SortMapLayer) (err error) {

	for _, s := range params {
		n.adaptors.MapLayer.Sort(&m.MapLayer{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

func (n *MapLayerEndpoint) Delete(mId int64) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	if _, err = n.adaptors.MapLayer.GetById(mId); err != nil {
		return
	}

	err = n.adaptors.MapLayer.Delete(mId)

	return
}

func (n *MapLayerEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.MapLayer, total int64, err error) {

	result, total, err = n.adaptors.MapLayer.List(limit, offset, order, sortBy)

	return
}
