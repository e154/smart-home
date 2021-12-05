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
	"fmt"
	"github.com/e154/smart-home/adaptors"
	"github.com/e154/smart-home/api/dto"
	"github.com/e154/smart-home/common"
	"github.com/e154/smart-home/endpoint"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/access_list"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

var (
	log = common.MustGetLogger("controllers")
)

// ControllerCommon ...
type ControllerCommon struct {
	adaptors   *adaptors.Adaptors
	accessList access_list.AccessListService
	endpoint   *endpoint.Endpoint
	dto        dto.Dto
}

// NewControllerCommon ...
func NewControllerCommon(adaptors *adaptors.Adaptors,
	accessList access_list.AccessListService,
	endpoint *endpoint.Endpoint) *ControllerCommon {
	return &ControllerCommon{
		dto:        dto.NewDto(),
		adaptors:   adaptors,
		accessList: accessList,
		endpoint:   endpoint,
	}
}

func (c ControllerCommon) currentUser(ctx context.Context) (*m.User, error) {

	user, ok := ctx.Value("currentUser").(*m.User)
	if !ok {
		return nil, fmt.Errorf("bad user object")
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
	json.NewEncoder(w).Encode(p)
}

func (c ControllerCommon) error(ctx context.Context, errs validator.ValidationErrorsTranslations, err error) error {
	if len(errs) > 0 {
		return c.prepareErrors(errs)
	}

	switch {
	case errors.Is(err, common.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (c ControllerCommon) Pagination(limit, offset uint32, order, sortBy string) (pagination common.PageParams) {

	pagination = common.PageParams{
		Limit:  200,
		Offset: 0,
		Order:  "desc",
		SortBy: "created_at",
	}

	if limit != 0 {
		pagination.Limit = int64(limit)
	}
	if offset != 0 {
		pagination.Offset = int64(offset)
	}
	if order != "" {
		pagination.Order = order
	}
	if sortBy != "" {
		pagination.SortBy = sortBy
	}

	return
}
