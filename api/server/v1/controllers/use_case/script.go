package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"errors"
)

func AddScript(script *m.Script, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation
	ok, errs = script.Valid()
	if len(errs) > 0 {
		return
	}

	if id, err = adaptors.Script.Add(script); err != nil {
		return
	}

	script.Id = id

	return
}

func GetScriptById(scriptId int64, adaptors *adaptors.Adaptors) (script *m.Script, err error) {

	script, err = adaptors.Script.GetById(scriptId)

	return
}

func UpdateScript(script *m.Script, adaptors *adaptors.Adaptors, core *core.Core) (ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = script.Valid()
	if len(errs) > 0 {
		return
	}

	if err = adaptors.Script.Update(script); err != nil {
		return
	}

	return
}

func GetScriptList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*m.Script, total int64, err error) {

	items, total, err = adaptors.Script.List(limit, offset, order, sortBy)

	return
}

func DeleteScriptById(scriptId int64, adaptors *adaptors.Adaptors, core *core.Core) (err error) {

	if scriptId == 0 {
		err = errors.New("script id is null")
		return
	}

	var script *m.Script
	if script, err = adaptors.Script.GetById(scriptId); err != nil {
		return
	}

	err = adaptors.Script.Delete(script.Id)

	return
}

