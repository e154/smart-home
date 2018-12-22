package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	m "github.com/e154/smart-home/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
	"strconv"
)

type ControllerRole struct {
	*ControllerCommon
}

func NewControllerRole(common *ControllerCommon) *ControllerRole {
	return &ControllerRole{ControllerCommon: common}
}

// Role godoc
// @tags role
// @Summary Add new role
// @Description
// @Produce json
// @Accept  json
// @Param role body models.NewRole true "role params"
// @Success 200 {object} models.NewObjectSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /role [post]
func (c ControllerRole) AddRole(ctx *gin.Context) {

	var params models.NewRole
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	role := &m.Role{
		Name:        params.Name,
		Description: params.Description,
	}

	_, errs, err := AddRole(role, c.adaptors)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("role", role).Send(ctx)
}

// Role godoc
// @tags role
// @Summary Show role
// @Description Get role by name
// @Produce json
// @Accept  json
// @Param name path string true "Role name"
// @Success 200 {object} models.ResponseRole
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /role/{name} [Get]
func (c ControllerRole) GetRoleByName(ctx *gin.Context) {

	name := ctx.Param("name")
	role, err := GetRoleByName(name, c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("role", role).Send(ctx)
}

// Role godoc
// @tags role
// @Summary get role access list
// @Description Get access list
// @Produce json
// @Accept  json
// @Param name path string true "Role name"
// @Success 200 {object} models.ResponseRole
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /role/{name}/access_list [Get]
func (c ControllerRole) GetAccessList(ctx *gin.Context) {

	name := ctx.Param("name")
	accessList, err := GetAccessList(name, c.adaptors, c.accessList)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("access_list", accessList).Send(ctx)
}

// Role godoc
// @tags role
// @Summary update role access list
// @Description Update access list
// @Produce json
// @Accept  json
// @Param name path string true "Role name"
// @Success 200 {object} models.ResponseRole
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /role/{name}/access_list [Put]
func (c ControllerRole) UpdateAccessList(ctx *gin.Context) {

	name := ctx.Param("name")
	role, err := GetRoleByName(name, c.adaptors)
	if err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("role", role).Send(ctx)
}

// Role godoc
// @tags role
// @Summary Update role
// @Description Update role by name
// @Produce json
// @Accept  json
// @Param  name path string true "Role name"
// @Param  role body models.UpdateRole true "Update role"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /role/{name} [Put]
func (c ControllerRole) UpdateRole(ctx *gin.Context) {

	name := ctx.Param("name")
	role, err := c.adaptors.Role.GetByName(name)
	if err != nil {
		NewError(404, err).Send(ctx)
		return
	}

	if err := ctx.ShouldBindJSON(&role); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	role.Name = name

	_, errs, err := UpdateRole(role, c.adaptors)
	if len(errs) > 0 {
		err400 := NewError(400)
		err400.ValidationToErrors(errs).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// Role godoc
// @tags role
// @Summary Role list
// @Description Get role list
// @Produce json
// @Accept  json
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Param order query string false "order" default(DESC)
// @Param sort_by query string false "sort_by" default(name)
// @Success 200 {object} models.ResponseRoleList
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /roles [Get]
func (c ControllerRole) GetRoleList(ctx *gin.Context) {

	var limit = 15
	var offset = 0
	var order = "DESC"
	var sortBy = "name"

	var err error
	if limit, err = strconv.Atoi(c.query(ctx, "limit")); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if offset, err = strconv.Atoi(c.query(ctx, "offset")); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	if c.query(ctx, "order") != "" {
		order = c.query(ctx, "order")
	}

	if c.query(ctx, "sort_by") != "" {
		sortBy = c.query(ctx, "sort_by")
	}

	items, total, err := GetRoleList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, int(total), items).Send(ctx)
	return
}

// Role godoc
// @tags role
// @Summary Delete role
// @Description Delete role by name
// @Produce json
// @Accept  json
// @Param  name path string true "Role name"
// @Success 200 {object} models.ResponseSuccess
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /role/{name} [Delete]
func (c ControllerRole) DeleteRoleByName(ctx *gin.Context) {

	name := ctx.Param("name")
	_, err := c.adaptors.Role.GetByName(name)
	if err != nil {
		NewError(404, err).Send(ctx)
		return
	}

	if err := DeleteRoleByName(name, c.adaptors); err != nil {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Send(ctx)
}

// Role godoc
// @tags role
// @Summary Search role
// @Description Search role by name
// @Produce json
// @Accept  json
// @Param query query string false "query"
// @Param limit query int true "limit" default(10)
// @Param offset query int true "offset" default(0)
// @Success 200 {object} models.SearchRoleResponse
// @Failure 400 {object} models.ErrorModel "some error"
// @Failure 404 {object} models.ErrorModel "some error"
// @Failure 500 {object} models.ErrorModel "some error"
// @Router /roles/search [Get]
func (c ControllerRole) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	roles, total, err := Search(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(map[string]interface{}{
		"roles": roles,
		"total": total,
	})
	resp.Send(ctx)
}
