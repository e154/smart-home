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

package web

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/e154/smart-home/common/apperr"
	"github.com/e154/smart-home/common/logger"
)

var (
	log = logger.MustGetLogger("web")
)

type crawler struct {
	headers            []map[string]string
	digestAuth         bool
	basicAuth          bool
	username, password string
	dig                *DigestHeaders
}

func New() Crawler {
	return &crawler{}
}

func (c *crawler) BasicAuth(username, password string) Crawler {
	c.username = username
	c.password = password
	c.basicAuth = true
	return c
}

func (c *crawler) DigestAuth(username, password string) Crawler {
	c.username = username
	c.password = password
	c.digestAuth = true
	return c
}

// Probe ...
func (c *crawler) Probe(options Request) (status int, body []byte, err error) {

	var req *http.Request
	if req, err = c.prepareRequest(options); err != nil {
		return
	}

	if err = c.auth(req, options.Url); err != nil {
		return
	}

	var timeout = time.Second * 2
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	var resp *http.Response
	if resp, err = c.doIt(req, timeout); err != nil {
		return
	}

	status = resp.StatusCode

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)

	return
}

// Probe ...
func (c *crawler) Download(options Request) (filePath string, err error) {

	defer func() {
		if err != nil {
			fmt.Printf("%v", err.Error())
		}
	}()

	var req *http.Request
	if req, err = c.prepareRequest(options); err != nil {
		return
	}

	if err = c.auth(req, options.Url); err != nil {
		return
	}

	var timeout = time.Second * 2
	if options.Timeout > 0 {
		timeout = options.Timeout
	}

	var resp *http.Response
	if resp, err = c.doIt(req, timeout); err != nil {
		return
	}

	defer resp.Body.Close()

	var f *os.File
	f, err = os.CreateTemp("/tmp", "*.jpeg")
	if err != nil {
		return
	}

	filePath = f.Name()
	if _, err = io.Copy(f, resp.Body); err != nil {
		return
	}

	log.Infof("file downloaded to '%s' ...", filePath)

	return
}

func (c *crawler) doIt(req *http.Request, timeout time.Duration) (resp *http.Response, err error) {

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

	resp, err = client.Do(req)

	return
}

func (c *crawler) prepareRequest(options Request) (req *http.Request, err error) {

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

	return
}

func (c *crawler) auth(req *http.Request, uri string) (err error) {

	if c.basicAuth {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
		c.headers = append(c.headers, map[string]string{
			"Authorization": "Basic " + auth,
		})
		req.Header.Add("Authorization", "Basic "+auth)
	}

	if c.digestAuth || c.dig != nil {
		c.dig = &DigestHeaders{}
		c.dig, err = c.dig.Auth(c.username, c.password, uri)
		if err != nil {
			return
		}
		c.dig.ApplyAuth(req)
	}
	return err
}
