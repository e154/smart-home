package use_case

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/system/core"
	"github.com/e154/smart-home/system/validation"
	"errors"
	"github.com/e154/smart-home/system/scripts"
	"github.com/e154/smart-home/api/server/v1/models"
	"github.com/jinzhu/copier"
)

func AddScript(script *m.Script, adaptors *adaptors.Adaptors, core *core.Core, scriptService *scripts.ScriptService) (ok bool, id int64, errs []*validation.Error, err error) {

	// validation
	ok, errs = script.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	var engine *scripts.Engine
	if engine, err = scriptService.NewEngine(script); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	if id, err = adaptors.Script.Add(script); err != nil {
		return
	}

	script.Id = id

	return
}

func GetScriptById(scriptId int64, adaptors *adaptors.Adaptors) (script *models.Script, err error) {

	var mScript *m.Script
	if mScript, err = adaptors.Script.GetById(scriptId); err != nil {
		return
	}

	script = &models.Script{}
	err = copier.Copy(&script, &mScript)

	return
}

func UpdateScript(script *m.Script, adaptors *adaptors.Adaptors, core *core.Core, scriptService *scripts.ScriptService) (ok bool, errs []*validation.Error, err error) {

	// validation
	ok, errs = script.Valid()
	if len(errs) > 0 || !ok {
		return
	}

	var engine *scripts.Engine
	if engine, err = scriptService.NewEngine(script); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	if err = adaptors.Script.Update(script); err != nil {
		return
	}

	return
}

func GetScriptList(limit, offset int64, order, sortBy string, adaptors *adaptors.Adaptors) (items []*models.Script, total int64, err error) {

	var mScripts []*m.Script
	if mScripts, total, err = adaptors.Script.List(limit, offset, order, sortBy); err != nil {
		return
	}

	items = make([]*models.Script, 0)

	for _, script := range mScripts {
		item := &models.Script{}
		if err = copier.Copy(&item, &script); err != nil {
			log.Error(err.Error())
			continue
		}
		items = append(items, item)
	}

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

func ExecuteScript(scriptId int64, adaptors *adaptors.Adaptors, core *core.Core, scriptService *scripts.ScriptService) (result string, err error) {

	var script *m.Script
	if script, err = adaptors.Script.GetById(scriptId); err != nil {
		return
	}

	var engine *scripts.Engine
	if engine, err = scriptService.NewEngine(script); err != nil {
		return
	}

	result, err = engine.DoFull()

	return
}

func SearchScript(query string, limit, offset int, adaptors *adaptors.Adaptors) (devices []*m.Script, total int64, err error) {

	devices, total, err = adaptors.Script.Search(query, limit, offset)

	return
}
