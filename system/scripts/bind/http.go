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

package bind

import (
	"github.com/e154/smart-home/common/web"
)

// HttpResponse ...
type HttpResponse struct {
	Body         string `json:"body"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

// HttpBind ...
type HttpBind struct {
	crawler web.Crawler
}

func NewHttpBind(crawler web.Crawler) *HttpBind {
	return &HttpBind{crawler: crawler}
}

// Get ...
func (h *HttpBind) Get(url string) (response HttpResponse) {
	//log.Infof("call [GET ] request %s", url)
	_, body, err := h.crawler.Probe(web.Request{Method: "GET", Url: url})
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return
	}
	response.Body = string(body)
	return
}

// Post ...
func (h *HttpBind) Post(url, data string) (response HttpResponse) {
	//log.Infof("call [POST] request %s", url)
	_, body, err := h.crawler.Probe(web.Request{Method: "POST", Url: url, Body: []byte(data)})
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return
	}
	response.Body = string(body)
	return
}
