package use_case

import (
	"errors"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/common/debug"
	"fmt"
)

type MapElementCommand struct {
	*CommonCommand
}

func NewMapElementCommand(common *CommonCommand) *MapElementCommand {
	return &MapElementCommand{
		CommonCommand: common,
	}
}

func (n *MapElementCommand) Add(params *m.MapElement) (result *m.MapElement, errs []*validation.Error, err error) {

	// validation
	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var id int64
	if id, err = n.adaptors.MapElement.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.MapElement.GetById(id)

	return
}

func (n *MapElementCommand) GetById(mId int64) (result *m.MapElement, err error) {

	result, err = n.adaptors.MapElement.GetById(mId)

	return
}

func (n *MapElementCommand) Update(params *m.MapElement) (result *m.MapElement, errs []*validation.Error, err error) {

	var m *m.MapElement
	if m, err = n.adaptors.MapElement.GetById(params.Id); err != nil {
		return
	}

	common.Copy(&m, &params, common.JsonEngine)

	// validation
	_, errs = m.Valid()
	if len(errs) > 0 {
		return
	}

	debug.Println(m)
	fmt.Println("----")
	if err = n.adaptors.MapElement.Update(m); err != nil {
		return
	}

	result, err = n.adaptors.MapElement.GetById(m.Id)

	return
}

func (n *MapElementCommand) Sort(params []*m.SortMapElement) (err error) {

	for _, s := range params {
		n.adaptors.MapElement.Sort(&m.MapElement{
			Id:     s.Id,
			Weight: s.Weight,
		})
	}

	return
}

func (n *MapElementCommand) Delete(mId int64) (err error) {

	if mId == 0 {
		err = errors.New("m id is null")
		return
	}

	if _, err = n.adaptors.MapElement.GetById(mId); err != nil {
		return
	}

	err = n.adaptors.MapElement.Delete(mId)

	return
}

func (n *MapElementCommand) GetList(limit, offset int64, order, sortBy string) (result []*m.MapElement, total int64, err error) {

	result, total, err = n.adaptors.MapElement.List(limit, offset, order, sortBy)

	return
}