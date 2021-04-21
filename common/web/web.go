// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2020, Filippov Alex
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
	"io/ioutil"
	"net/http"
)

func Crawler(uri string, headers ...map[string]string) (body []byte, err error) {

	var req *http.Request
	req, err = http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	} else {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15")
	}

	client := &http.Client{}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	return
}
