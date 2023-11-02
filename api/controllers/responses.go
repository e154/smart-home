package controllers

import (
	"github.com/labstack/echo/v4"

	"github.com/e154/smart-home/common"
)

// Pagination ...
type Pagination struct {
	Page  int64  `json:"page"`
	Total int64  `json:"total"`
	Limit uint64 `json:"limit"`
}

// GenericMeta ...
type GenericMeta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
	Sort       *string     `json:"sort,omitempty"`
}

func NewGenericMeta(ctx echo.Context) *GenericMeta {
	return &GenericMeta{}
}

func (g *GenericMeta) WithPagination(pagination Pagination) *GenericMeta {
	g.Pagination = &pagination
	return g
}

func (g *GenericMeta) WithSort(sort string) *GenericMeta {
	g.Sort = common.String(sort)
	return g
}

// ErrorField ...
type ErrorField struct {
	Name    *string `json:"name,omitempty"`
	Message *string `json:"message,omitempty"`
}

// ErrorBase ...
type ErrorBase struct {
	Code    *string      `json:"code,omitempty"`
	Message *string      `json:"message,omitempty"`
	Fields  []ErrorField `json:"fields,omitempty"`
}

// Error ...
type Error struct {
	Error *ErrorBase   `json:"error,omitempty"`
	Meta  *GenericMeta `json:"meta,omitempty"`
}

// Data ...
type Data struct {
	ID interface{} `json:"id,omitempty"`
}

// Success ...
type Success struct {
	Data  interface{}  `json:"data,omitempty"`
	Items interface{}  `json:"items"`
	Meta  *GenericMeta `json:"meta,omitempty"`
}

// ResponseWithObj ...
func ResponseWithObj(ctx echo.Context, obj interface{}) interface{} {
	return obj
}

// ResponseWithList ...
func ResponseWithList(ctx echo.Context, items interface{}, total int64, pagination common.PageParams) *Success {
	return &Success{
		Items: items,
		Meta: NewGenericMeta(ctx).
			WithPagination(Pagination{
				Page:  pagination.PageReq,
				Total: total,
				Limit: uint64(pagination.Limit),
			}).
			WithSort(pagination.SortReq),
	}
}

func ResponseWithError(ctx echo.Context, err *ErrorBase) *Error {
	return &Error{
		Error: err,
	}
}
