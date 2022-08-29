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
	"net"
	"net/http"
	"time"

	"github.com/e154/smart-home/common/apperr"
)

type crawler struct {
}

func New() Crawler {
	return &crawler{}
}

// Probe ...
func (crawler) Probe(options Request) (status int, body []byte, err error) {
	return Probe(options)
}

// Probe ...
func Probe(options Request) (status int, body []byte, err error) {

	if options.Url == "" {
		err = apperr.ErrBadRequestParams
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

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15")

	// set headers
	for _, values := range options.Headers {
		for k, v := range values {
			req.Header.Set(k, v)
		}
	}

	var timeout = time.Second * 2
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	netTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
		TLSHandshakeTimeout: timeout,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: netTransport,
	}

	var resp *http.Response
	if resp, err = client.Do(req); err != nil {
		return
	}

	status = resp.StatusCode

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	return
}
