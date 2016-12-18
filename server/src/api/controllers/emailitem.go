package controllers

import (
	"fmt"
	"../models"
)

type EmailItemController struct {
	CommonController
}

func (i *EmailItemController) Post() {

	_, b, valid, err := models.EmailItemAddNew(i.Ctx.Input.RequestBody)
	if err != nil {
		i.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf( "%s: %s\r", err.Key, err.Message)
		}
		i.ErrHan(403, err.Error())
		return
	}

}

func (i *EmailItemController) GetOne() {

	name := i.Ctx.Input.Param(":name")

	item, err := models.EmailItemGet(name)
	if err != nil {
		i.ErrHan(403, err.Error())
		return
	}

	i.Data["json"] = &map[string]interface {}{"status": "success", "item": item}
	i.ServeJSON()
}

func (i *EmailItemController) GetAll() {

	count, items, err := models.EmailItemGetSortedList()
	if err != nil {
		i.ErrHan(403, err.Error())
		return
	}

	i.Data["json"] = &map[string]interface {}{"status": "success", "items": items, "amount": count}
	i.ServeJSON()
}

func (i *EmailItemController) Put() {

	name := i.Ctx.Input.Param(":name")

	_, b, valid, err := models.EmailItemUpdate(i.Ctx.Input.RequestBody, name)
	if err != nil {
		i.ErrHan(403, err.Error())
		return
	}

	if !b {
		var msg string
		for _, err := range valid.Errors {
			msg += fmt.Sprintf( "%s: %s\r", err.Key, err.Message)
		}
		i.ErrHan(403, err.Error())
		return
	}

}

func (i *EmailItemController) Delete() {

	name := i.Ctx.Input.Param(":name")
	if err := models.EmailItemDelete(name); err != nil {
		i.ErrHan(403, err.Error())
		return
	}

}

func (i *EmailItemController) GetTree() {

	tree, err := models.EmailItemGetTree()
	if err != nil {
		i.ErrHan(403, err.Error())
		return
	}

	i.Data["json"] = &map[string]interface {}{"status": "success", "tree": tree}
	i.ServeJSON()
}

func (i *EmailItemController) UpdateTree() {

	if err := models.EmailItemUpdateTree(i.Ctx.Input.RequestBody); err != nil {
		i.ErrHan(403, err.Error())
		return
	}
}
