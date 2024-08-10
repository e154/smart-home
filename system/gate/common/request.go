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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// HTTPRequest is a serializable version of http.Request ( with only usefull fields )
type HTTPRequest struct {
	Method        string
	URL           string
	Header        map[string][]string
	ContentLength int64
	WS            bool
	Body          []byte
}

// SerializeHTTPRequest create a new HTTPRequest from a http.Request
func SerializeHTTPRequest(req *http.Request) (r *HTTPRequest) {
	body, _ := io.ReadAll(req.Body)
	r = &HTTPRequest{
		URL:           req.URL.String(),
		Method:        req.Method,
		Header:        req.Header,
		ContentLength: req.ContentLength,
		Body:          body,
	}
	return
}

// DeserializeHTTPRequest create a new http.Request from a HTTPRequest
func DeserializeHTTPRequest(jsonRequest []byte) (r *http.Request, err error) {

	req := &HTTPRequest{}
	err = json.Unmarshal(jsonRequest, req)
	if err != nil {
		err = fmt.Errorf("unable to deserialize json http request : %s", err)
		return
	}

	var uri *url.URL
	if uri, err = url.Parse(req.URL); err != nil {
		err = fmt.Errorf("unable to parse url : %s", err.Error())
		return
	}

	r = &http.Request{
		Method:        req.Method,
		URL:           uri,
		Header:        req.Header,
		ContentLength: req.ContentLength,
		Body:          io.NopCloser(bytes.NewReader(req.Body)),
	}
	return
}
