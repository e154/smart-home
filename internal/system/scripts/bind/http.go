// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2023, Filippov Alex
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
	web2 "github.com/e154/smart-home/internal/common/web"
	"github.com/e154/smart-home/pkg/web"
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
	headers []map[string]string
}

func NewHttpBind() *HttpBind {
	return &HttpBind{crawler: web2.New()}
}

// Get ...
func (h *HttpBind) Get(url string) (response HttpResponse) {
	//log.Infof("call [GET ] request %s", url)
	_, body, err := h.crawler.Probe(web.Request{Method: "GET", Url: url, Headers: h.headers})
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
	_, body, err := h.crawler.Probe(web.Request{Method: "POST", Url: url, Headers: h.headers, Body: []byte(data)})
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return
	}
	response.Body = string(body)
	return
}

// Put ...
func (h *HttpBind) Put(url, data string) (response HttpResponse) {
	//log.Infof("call [PUT] request %s", url)
	_, body, err := h.crawler.Probe(web.Request{Method: "PUT", Url: url, Headers: h.headers, Body: []byte(data)})
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return
	}
	response.Body = string(body)
	return
}

// Delete ...
func (h *HttpBind) Delete(url string) (response HttpResponse) {
	//log.Infof("call [DELETE] request %s", url)
	_, body, err := h.crawler.Probe(web.Request{Method: "DELETE", Url: url, Headers: h.headers})
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return
	}
	response.Body = string(body)
	return
}

func (h *HttpBind) Headers(headers []map[string]string) *HttpBind {
	return &HttpBind{
		crawler: h.crawler,
		headers: headers,
	}
}

func (h *HttpBind) BasicAuth(username, password string) *HttpBind {
	return &HttpBind{crawler: web2.New().BasicAuth(username, password)}
}

func (h *HttpBind) DigestAuth(username, password string) *HttpBind {
	return &HttpBind{crawler: web2.New().DigestAuth(username, password)}
}

func (h *HttpBind) Download(uri string) (response HttpResponse) {
	filePath, err := h.crawler.Download(web.Request{Method: "GET", Url: uri})
	if err != nil {
		response.Error = true
		response.ErrorMessage = err.Error()
		return
	}
	response.Body = filePath
	return
}
