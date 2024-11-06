// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2024, Filippov Alex
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

package rbac_http

import (
	"context"
	"net/http"
	"strings"

	"github.com/e154/smart-home/pkg/apperr"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

var _ plugins.HttpAccessFilter = (*HttpAccessFilter)(nil)

type HttpAccessFilter struct {
	config        *m.AppConfig
	Authenticator plugins.Authorization
}

func NewHttpAccessFilter(config *m.AppConfig, authenticator plugins.Authorization) plugins.HttpAccessFilter {
	return &HttpAccessFilter{
		config:        config,
		Authenticator: authenticator,
	}
}

func (f *HttpAccessFilter) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestURI := r.URL
		method := strings.ToLower(r.Method)

		var accessToken = f.getAccessToken(r)
		if accessToken == "" {
			f.HTTP401(w, apperr.ErrUnauthorized)
			return
		}

		user, root, err := f.Authenticator.AuthREST(r.Context(), accessToken, requestURI, method)
		if err != nil {
			f.HTTP401(w, apperr.ErrUnauthorized)
			return
		}
		if user == nil {
			f.HTTP401(w, apperr.ErrUnauthorized)
		}

		ctx := context.WithValue(r.Context(), "currentUser", user)
		ctx = context.WithValue(r.Context(), "root", f.config.RootMode || root)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (f *HttpAccessFilter) getAccessToken(r *http.Request) (accessToken string) {
	accessToken = r.Header.Get("authorization")
	if accessToken != "" {
		return
	}
	accessToken = r.URL.Query().Get("access_token")
	return
}

// HTTP401 ...
func (f *HttpAccessFilter) HTTP401(w http.ResponseWriter, err error) {
	http.Error(w, "UNAUTHORIZED", http.StatusUnauthorized)
}
