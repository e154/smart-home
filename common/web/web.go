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

package web

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/e154/smart-home/common"
)

// Request ...
type Request struct {
	Method  string
	Url     string
	Body    []byte
	headers []map[string]string
	Timeout time.Duration
}

// Crawler ...
func Crawler(options Request) (body []byte, err error) {

	if options.Url == "" {
		err = common.ErrBadRequestParams
		return
	}

	var r io.Reader = nil
	switch options.Method {
	case "POST", "PUT", "UPDATE", "PATCH":
		if options.Body != nil && len(options.Body) > 0 {
			r = bytes.NewReader(options.Body)
		}
	}

	var req *http.Request
	if req, err = http.NewRequest(options.Method, options.Url, r); err != nil {
		return
	}

	if len(options.headers) > 0 {
		for k, v := range options.headers[0] {
			req.Header.Set(k, v)
		}
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15")
	}

	var timeout = time.Second * 2
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	client := &http.Client{
		Timeout: timeout,
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	return
}
