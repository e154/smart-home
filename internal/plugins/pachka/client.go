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

package pachka

import (
	"encoding/json"
	"fmt"
	"io"

	web2 "github.com/e154/smart-home/internal/system/web"
	"github.com/e154/smart-home/pkg/web"
)

const (
	baseUrl = "https://api.pachca.com/api/shared/v1"
)

type Client struct {
	accessToken string
}

func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
	}
}

func (p *Client) SendMsg(content string, entityId int64) (responseMessage *ResponseMessage, errs []*ErrorItem, err error) {

	j, _ := json.Marshal(map[string]interface{}{
		"message": RequestMessage{
			EntityId: entityId,
			Content:  content,
		},
	})

	crawler := web2.New()
	request := web.Request{
		Method: "POST",
		Url:    fmt.Sprintf("%s/messages", baseUrl),
		Body:   j,
		Headers: []map[string]string{
			{"Authorization": "Bearer " + p.accessToken},
			{"Content-Type": "application/json; charset=utf-8"},
		},
	}

	var code int
	var body []byte
	code, body, err = crawler.Probe(request)
	if err != nil {
		return
	}

	switch code {
	case 200, 201:
	default:
		_err := &ResponseError{}
		err = json.Unmarshal(body, &_err)
		errs = _err.Errors
		return
	}

	responseMessage = &ResponseMessage{}
	err = json.Unmarshal(body, &responseMessage)
	if err == io.EOF {
		err = nil
	}

	return
}
