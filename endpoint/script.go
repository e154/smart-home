package endpoint

import (
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/validation"
	"github.com/e154/smart-home/system/scripts"
	m "github.com/e154/smart-home/models"
	"errors"
)

type ScriptEndpoint struct {
	*CommonEndpoint
}

func NewScriptEndpoint(common *CommonEndpoint) *ScriptEndpoint {
	return &ScriptEndpoint{
		CommonEndpoint: common,
	}
}

func (n *ScriptEndpoint) Add(params *m.Script) (result *m.Script, errs []*validation.Error, err error) {

	_, errs = params.Valid()
	if len(errs) > 0 {
		return
	}

	var engine *scripts.Engine
	if engine, err = n.scriptService.NewEngine(params); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	var id int64
	if id, err = n.adaptors.Script.Add(params); err != nil {
		return
	}

	result, err = n.adaptors.Script.GetById(id)

	return
}

func (n *ScriptEndpoint) GetById(scriptId int64) (result *m.Script, err error) {

	result, err = n.adaptors.Script.GetById(scriptId)

	return
}

func (n *ScriptEndpoint) Update(params *m.Script) (result *m.Script, errs []*validation.Error, err error) {

	var script *m.Script
	if script, err = n.adaptors.Script.GetById(params.Id); err != nil {
		return
	}

	if err = common.Copy(&script, &params); err != nil {
		return
	}

	// validation
	_, errs = script.Valid()
	if len(errs) > 0 {
		return
	}

	var engine *scripts.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	if err = n.adaptors.Script.Update(script); err != nil {
		return
	}

	result, err = n.adaptors.Script.GetById(script.Id)

	return
}

func (n *ScriptEndpoint) GetList(limit, offset int64, order, sortBy string) (result []*m.Script, total int64, err error) {

	result, total, err = n.adaptors.Script.List(limit, offset, order, sortBy)

	return
}

func (n *ScriptEndpoint) DeleteScriptById(scriptId int64) (err error) {

	if scriptId == 0 {
		err = errors.New("script id is null")
		return
	}

	var script *m.Script
	if script, err = n.adaptors.Script.GetById(scriptId); err != nil {
		return
	}

	err = n.adaptors.Script.Delete(script.Id)

	return
}

func (n *ScriptEndpoint) Execute(scriptId int64) (result string, err error) {

	var script *m.Script
	if script, err = n.adaptors.Script.GetById(scriptId); err != nil {
		return
	}

	var engine *scripts.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		return
	}

	result, err = engine.DoFull()

	return
}

func (n *ScriptEndpoint) ExecuteSource(script *m.Script) (result string, err error) {

	var engine *scripts.Engine
	if engine, err = n.scriptService.NewEngine(script); err != nil {
		return
	}

	if err = engine.Compile(); err != nil {
		return
	}

	result, err = engine.DoFull()

	return
}

func (n *ScriptEndpoint) Search(query string, limit, offset int) (devices []*m.Script, total int64, err error) {

	devices, total, err = n.adaptors.Script.Search(query, limit, offset)

	return
}
