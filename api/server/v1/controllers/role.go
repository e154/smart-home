package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/e154/smart-home/api/server/v1/models"
	. "github.com/e154/smart-home/api/server/v1/controllers/use_case"
)

type ControllerRole struct {
	*ControllerCommon
}

func NewControllerRole(common *ControllerCommon) *ControllerRole {
	return &ControllerRole{ControllerCommon: common}
}

// swagger:operation POST /role roleAdd
// ---
// parameters:
// - description: role params
//   in: body
//   name: role
//   required: true
//   schema:
//     $ref: '#/definitions/NewRole'
//     type: object
// summary: add new role
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Role'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) Add(ctx *gin.Context) {

	var params models.NewRole
	if err := ctx.ShouldBindJSON(&params); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	role, errs, err := AddRole(params, c.adaptors)
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
	resp.SetData(role).Send(ctx)
}

// swagger:operation GET /role/{name} roleGetById
// ---
// parameters:
// - description: Role name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get role by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Role'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) GetByName(ctx *gin.Context) {

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
	resp.SetData(role).Send(ctx)
}

// swagger:operation GET /role/{name}/access_list roleGetById
// ---
// parameters:
// - description: Role name
//   in: path
//   name: name
//   required: true
//   type: string
// summary: get access list by role name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       type: object
//       properties:
//         access_list:
//           $ref: '#/definitions/AccessList'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
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

// swagger:operation PUT /role/{name}/access_list roleUpdateById
// ---
// parameters:
// - description: Role name
//   in: path
//   name: name
//   required: true
//   type: string
// - description: Update access list params
//   in: body
//   name: access_list_diff
//   required: true
//   schema:
//     $ref: '#/definitions/AccessListDiff'
// summary: update role access list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) UpdateAccessList(ctx *gin.Context) {

	accessListDif := make(map[string]map[string]bool)
	if err := ctx.ShouldBindJSON(&accessListDif); err != nil {
		log.Error(err.Error())
		NewError(400, err).Send(ctx)
		return
	}

	name := ctx.Param("name")
	if err := UpdateAccessList(name, accessListDif, c.adaptors); err != nil {
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

// swagger:operation PUT /role/{name} roleUpdateById
// ---
// parameters:
// - description: Role ID
//   in: path
//   name: name
//   required: true
//   type: string
// - description: Update role params
//   in: body
//   name: role
//   required: true
//   schema:
//     $ref: '#/definitions/UpdateRole'
// summary: update role by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//     description: OK
//     schema:
//       $ref: '#/definitions/Role'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) Update(ctx *gin.Context) {

	name := ctx.Param("name")
	role := &models.UpdateRole{}
	if err := ctx.ShouldBindJSON(&role); err != nil {
		NewError(400, err).Send(ctx)
		return
	}

	role.Name = name

	result, errs, err := UpdateRole(role, c.adaptors)
	if len(errs) > 0 {
		code := 500
		if err.Error() == "record not found" {
			code = 404
		}
		NewError(code, err).Send(ctx)
		return
	}

	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.SetData(result).Send(ctx)
}

// swagger:operation GET /roles roleList
// ---
// summary: get role list
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// parameters:
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// - default: DESC
//   description: order
//   in: query
//   name: order
//   type: string
// - default: id
//   description: sort_by
//   in: query
//   name: sort_by
//   type: string
// responses:
//   "200":
//	   $ref: '#/responses/RoleList'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) GetList(ctx *gin.Context) {

	_, sortBy, order, limit, offset := c.list(ctx)
	items, total, err := GetRoleList(int64(limit), int64(offset), order, sortBy, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Page(limit, offset, total, items).Send(ctx)
	return
}

// swagger:operation DELETE /role/{name} roleDeleteById
// ---
// parameters:
// - description: Role ID
//   in: path
//   name: name
//   required: true
//   type: string
// summary: delete role by name
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// responses:
//   "200":
//	   $ref: '#/responses/Success'
//   "400":
//	   $ref: '#/responses/Error'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "404":
//	   $ref: '#/responses/Error'
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) Delete(ctx *gin.Context) {

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

// swagger:operation GET /roles/search roleSearch
// ---
// summary: search role
// description:
// security:
// - ApiKeyAuth: []
// tags:
// - role
// parameters:
// - description: query
//   in: query
//   name: query
//   type: string
// - default: 10
//   description: limit
//   in: query
//   name: limit
//   required: true
//   type: integer
// - default: 0
//   description: offset
//   in: query
//   name: offset
//   required: true
//   type: integer
// responses:
//   "200":
//	   $ref: '#/responses/RoleSearch'
//   "401":
//     description: "Unauthorized"
//   "403":
//     description: "Forbidden"
//   "500":
//	   $ref: '#/responses/Error'
func (c ControllerRole) Search(ctx *gin.Context) {

	query, limit, offset := c.select2(ctx)
	roles, _, err := SearchRole(query, limit, offset, c.adaptors)
	if err != nil {
		NewError(500, err).Send(ctx)
		return
	}

	resp := NewSuccess()
	resp.Item("roles", roles)
	resp.Send(ctx)
}
