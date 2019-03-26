package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/common"
	"errors"
)

type MapLayerCommand struct {
	*CommonCommand
}

func NewMapLayerCommand(common *CommonCommand) *MapLayerCommand {
	return &MapLayerCommand{
		CommonCommand: common,
	}
}

func (n *MapLayerCommand) Add(params *m.MapLayer) (result *m.MapLayer, errs []*validation.Error, err error) {

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

func (n *MapLayerCommand) GetById(mId int64) (result *m.MapLayer, err error) {

	result, err = n.adaptors.MapLayer.GetById(mId)

	return
}

func (n *MapLayerCommand) Update(params *m.MapLayer) (result *m.MapLayer, errs []*validation.Error, err error) {

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

func (n *MapLayerCommand) Sort(params []*m.SortMapLayer) (err error) {

	for _, s := range params {
		n.adaptors.MapLayer.Sort(&m.MapLayer{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

func (n *MapLayerCommand) Delete(mId int64) (err error) {

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

func (n *MapLayerCommand) GetList(limit, offset int64, order, sortBy string) (result []*m.MapLayer, total int64, err error) {

	result, total, err = n.adaptors.MapLayer.List(limit, offset, order, sortBy)

	return
}
