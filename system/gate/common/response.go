// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2023, Filippov Alex
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

package common

import (
	"fmt"
	"log"
	"net/http"
)

// HTTPResponse is a serializable version of http.Response ( with only useful fields )
type HTTPResponse struct {
	StatusCode    int
	Header        http.Header
	ContentLength int64
}

// SerializeHTTPResponse create a new HTTPResponse from a http.Response
func SerializeHTTPResponse(resp *http.Response) *HTTPResponse {
	r := new(HTTPResponse)
	r.StatusCode = resp.StatusCode
	r.Header = resp.Header
	r.ContentLength = resp.ContentLength
	return r
}

// NewHTTPResponse creates a new HTTPResponse
func NewHTTPResponse() (r *HTTPResponse) {
	r = new(HTTPResponse)
	r.Header = make(http.Header)
	return
}

// ProxyError log error and return a HTTP 526 error with the message
func ProxyError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), 526)
}

// ProxyErrorf log error and return a HTTP 526 error with the message
func ProxyErrorf(w http.ResponseWriter, format string, args ...interface{}) {
	ProxyError(w, fmt.Errorf(format, args...))
}
