// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2021, Filippov Alex
//
// This library is free software: you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 3 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Library General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library.  If not, see
// <https://www.gnu.org/licenses/>.

package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/dto"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/e154/smart-home/system/stream"
)

var (
	log = logger.MustGetLogger("controllers")
)

// ControllerCommon ...
type ControllerCommon struct {
	adaptors   *adaptors.Adaptors
	accessList access_list.AccessListService
	endpoint   *endpoint.Endpoint
	dto        dto.Dto
	stream     *stream.Stream
}

// NewControllerCommon ...
func NewControllerCommon(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	endpoint *endpoint.Endpoint,
	stream *stream.Stream) *ControllerCommon {
	return &ControllerCommon{
		dto:        dto.NewDto(),
		adaptors:   adaptors,
		accessList: accessList,
		endpoint:   endpoint,
		stream:     stream,
	}
}

func (c ControllerCommon) currentUser(ctx context.Context) (*m.User, error) {

	user, ok := ctx.Value("currentUser").(*m.User)
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

func (c ControllerCommon) prepareErrors(errs validator.ValidationErrorsTranslations) error {
	if len(errs) > 0 {
		st := status.New(codes.InvalidArgument, "One or more fields are invalid")
		for k, v := range errs {
			st, _ = st.WithDetails(&errdetails.BadRequest_FieldViolation{
				Field:       k,
				Description: v,
			})
		}
		return st.Err()
	}
	return nil
}

func (c ControllerCommon) writeErr(code int, body string, w http.ResponseWriter) {
	http.Error(w, body, code)
}

func (c ControllerCommon) writeSuccess(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (c ControllerCommon) writeJson(w http.ResponseWriter, p interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(p)
}

func (c ControllerCommon) error(_ context.Context, errs validator.ValidationErrorsTranslations, err error) error {
	if len(errs) > 0 {
		return c.prepareErrors(errs)
	}

	//defer func() {
	//	fmt.Println("-------")
	//	fmt.Println(errors.Cause(err).Error())
	//	fmt.Println("-------")
	//
	//	for {
	//		fmt.Println("--->", err.Error())
	//		err = errors.Unwrap(err)
	//		if err == nil {
	//			return
	//		}
	//	}
	//}()

	switch {
	case errors.Is(err, apperr.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, apperr.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, apperr.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, apperr.ErrUnimplemented):
		return status.Error(codes.Unimplemented, err.Error())
	default:
		log.Errorf("%+v\n", err)
		return status.Error(codes.Internal, err.Error())
	}
}

// Pagination ...
func (c ControllerCommon) Pagination(page, limit uint64, sort string) (pagination common.PageParams) {

	if sort == "" {
		sort = "-createdAt"
	}

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 200
	}

	pagination = common.PageParams{
		Limit:   int64(limit),
		Offset:  int64(page*limit - limit),
		Order:   "desc",
		SortBy:  "created_at",
		PageReq: page,
		SortReq: sort,
	}

	if len(sort) > 1 {
		firstChar := string([]rune(sort)[0])
		switch firstChar {
		case "+":
			pagination.Order = "asc"
		case "-":
			pagination.Order = "desc"
		default:
			//...
		}

		sort = strings.Replace(sort, firstChar, "", 1)
		pagination.SortBy = strcase.ToSnake(sort)
	}

	return
}

// Search ...
func (c ControllerCommon) Search(query string, limit, offset int64) (search common.SearchParams) {

	search = common.SearchParams{
		Query:  query,
		Limit:  200,
		Offset: 0,
	}

	if limit > 0 {
		search.Limit = limit
	}
	if offset > 0 {
		search.Offset = offset
	}

	return
}
