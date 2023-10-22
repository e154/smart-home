package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/e154/smart-home/api/dto"
	"github.com/e154/smart-home/system/access_list"

	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/validation"
)

var (
	log = logger.MustGetLogger("controllers")
)

// ControllerCommon ...
type ControllerCommon struct {
	endpoint   *endpoint.Endpoint
	accessList access_list.AccessListService
	validation *validation.Validate
	dto        dto.Dto
	appConfig  *m.AppConfig
}

// NewControllerCommon ...
func NewControllerCommon(endpoint *endpoint.Endpoint,
	accessList access_list.AccessListService,
	appConfig *m.AppConfig,
	validation *validation.Validate) *ControllerCommon {
	return &ControllerCommon{
		endpoint:   endpoint,
		appConfig:  appConfig,
		validation: validation,
		accessList: accessList,
		dto:        dto.NewDto(),
	}
}

func (c ControllerCommon) Body(ctx echo.Context, obj interface{}) error {
	dec := json.NewDecoder(ctx.Request().Body)
	if err := dec.Decode(obj); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			return apperr.ErrorWithCode("BAD_REQUEST", err.Error(), apperr.ErrUnknownField)
		}
		return apperr.ErrorWithCode("BAD_JSON_REQUEST", err.Error(), apperr.ErrBadJSONRequest)
	}
	return nil
}

// HTTP200 ...
func (c ControllerCommon) HTTP200(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusOK, data)
}

// HTTP201 ...
func (c ControllerCommon) HTTP201(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusCreated, data)
}

// HTTP401 ...
func (c ControllerCommon) HTTP401(ctx echo.Context, err error) error {
	e := apperr.GetError(err)
	if e != nil {
		return ctx.JSON(http.StatusUnauthorized, ResponseWithError(ctx, &ErrorBase{
			Code:    common.String(e.Code()),
			Message: common.String(e.Message()),
		}))
	}
	return ctx.JSON(http.StatusUnauthorized, ResponseWithError(ctx, &ErrorBase{
		Code: common.String("UNAUTHORIZED"),
	}))
}

// HTTP404 ...
func (c ControllerCommon) HTTP404(ctx echo.Context, err error) error {
	code := common.String("NOT_FOUND")
	message := common.String(err.Error())
	e := apperr.GetError(err)
	if e != nil {
		code = common.String(e.Code())
		message = common.String(e.Message())
	}
	return ctx.JSON(http.StatusNotFound, ResponseWithError(ctx, &ErrorBase{
		Code:    code,
		Message: message,
	}))
}

// HTTP400 ...
func (c ControllerCommon) HTTP400(ctx echo.Context, err error) error {
	code := common.String("BAD_REQUEST")
	message := common.String(err.Error())
	e := apperr.GetError(err)
	if e != nil {
		code = common.String(e.Code())
		message = common.String(e.Message())
	}
	return ctx.JSON(http.StatusBadRequest, ResponseWithError(ctx, &ErrorBase{
		Code:    code,
		Message: message,
	}))
}

// HTTP409 ...
func (c ControllerCommon) HTTP409(ctx echo.Context, err error) error {
	code := common.String("CONFLICT")
	message := common.String(err.Error())
	e := apperr.GetError(err)
	if e != nil {
		code = common.String(e.Code())
		message = common.String(e.Message())
	}
	return ctx.JSON(http.StatusConflict, ResponseWithError(ctx, &ErrorBase{
		Code:    code,
		Message: message,
	}))
}

// HTTP500 ...
func (c ControllerCommon) HTTP500(ctx echo.Context, err error) error {
	code := common.String("INTERNAL_ERROR")
	message := common.String(err.Error())
	e := apperr.GetError(err)
	if e != nil {
		code = common.String(e.Code())
		message = common.String(e.Message())
	}
	return ctx.JSON(http.StatusInternalServerError, ResponseWithError(ctx, &ErrorBase{
		Code:    code,
		Message: message,
	}))
}

// HTTP422 ...
func (c ControllerCommon) HTTP422(ctx echo.Context, err error) error {

	var fields []ErrorField

	respErr := ErrorBase{
		Code: common.String("UNPROCESSABLE_ERROR"),
	}

	e := apperr.GetError(err)
	if e != nil {
		errs := e.ValidationErrors()

		for fieldName, desc := range errs {
			// update field name
			fieldNameArr := strings.Split(fieldName, ".")
			fieldName = fieldNameArr[len(fieldNameArr)-1]

			fields = append(fields, ErrorField{
				Name:    common.String(fieldName),
				Message: common.String(desc),
			})
		}

		respErr.Code = common.String(e.Code())
		respErr.Message = common.String(e.Message())
		respErr.Fields = fields
	}

	return ctx.JSON(http.StatusUnprocessableEntity, ResponseWithError(ctx, &respErr))
}

// HTTP501 ...
func (c ControllerCommon) HTTP501(ctx echo.Context, data interface{}) error {
	return ctx.JSON(http.StatusNotImplemented, data)
}

// Pagination ...
func (c ControllerCommon) Pagination(page, limit *uint64, sort *string) (pagination common.PageParams) {

	pagination = common.PageParams{
		Limit:   200,
		Offset:  0,
		Order:   "desc",
		SortBy:  "created_at",
		PageReq: 1,
		SortReq: "-created_at",
	}

	if limit != nil {
		pagination.Limit = int64(*limit)
	}
	if page != nil {
		pagination.PageReq = int64(*page)
	}

	pagination.Offset = pagination.Limit * (pagination.PageReq - 1)
	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	if sort != nil && len(*sort) > 1 {
		pagination.SortReq = *sort
		firstChar := string([]rune(*sort)[0])
		switch firstChar {
		case "+":
			pagination.Order = "asc"
		case "-":
			pagination.Order = "desc"
		}

		// ToSnake converts a string to snake_case
		pagination.SortBy = strcase.ToSnake(strings.Replace(*sort, firstChar, "", 1))
	}

	return
}

// Search ...
func (c ControllerCommon) Search(query *string, limit, offset *int64) (search common.SearchParams) {

	search = common.SearchParams{
		Query:  common.StringValue(query),
		Limit:  200,
		Offset: 0,
	}

	if limit != nil {
		search.Limit = common.Int64Value(limit)
	}
	if offset != nil {
		search.Offset = common.Int64Value(offset)
	}

	return
}

// ERROR ...
func (c ControllerCommon) ERROR(ctx echo.Context, err error) error {
	switch {
	case errors.Is(err, apperr.ErrAlreadyExists):
		return c.HTTP409(ctx, err)
	case errors.Is(err, apperr.ErrInternal):
		return c.HTTP500(ctx, err)
	case errors.Is(err, apperr.ErrInvalidRequest):
		return c.HTTP422(ctx, err)
	case errors.Is(err, apperr.ErrNotFound):
		return c.HTTP404(ctx, err)
	case errors.Is(err, apperr.ErrAccessDenied):
		return c.HTTP401(ctx, err)
	case errors.Is(err, apperr.ErrUnknownField):
		return c.HTTP400(ctx, err)
	case errors.Is(err, apperr.ErrBadJSONRequest):
		return c.HTTP400(ctx, err)
	default:
		var bodyStr string
		body, _ := io.ReadAll(ctx.Request().Body)
		if len(body) > 0 {
			bodyStr = string(body)
		}
		url := ctx.Request().URL.String()
		log.Warnf("unknown err type %v for uri %s and body %q", err, url, bodyStr)
	}
	log.Error(err.Error())
	return nil
}

func (c ControllerCommon) currentUser(ctx echo.Context) (*m.User, error) {

	user, ok := ctx.Get("currentUser").(*m.User)
	if !ok {
		return nil, errors.Wrap(apperr.ErrBadRequestParams, "bad user object")
	}

	return user, nil
}

func (c ControllerCommon) parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	str, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(str)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}

	return cs[:s], cs[s+1:], true
}
