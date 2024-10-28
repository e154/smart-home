// This file is part of the Smart Home
// Program complex distribution https://github.com/e154/smart-home
// Copyright (C) 2016-2024, Filippov Alex
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
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/e154/smart-home/pkg/apperr"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/web"
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

func New() web.Crawler {
	return &crawler{}
}

func (c *crawler) BasicAuth(username, password string) web.Crawler {
	c.username = username
	c.password = password
	c.basicAuth = true
	return c
}

func (c *crawler) DigestAuth(username, password string) web.Crawler {
	c.username = username
	c.password = password
	c.digestAuth = true
	return c
}

// Probe ...
func (c *crawler) Probe(options web.Request) (status int, body []byte, err error) {

	var req *http.Request
	if req, err = c.prepareRequest(options); err != nil {
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

func (c *crawler) HEAD(options web.Request) (resp *http.Response, err error) {

	options.Method = "HEAD"

	// get heads
	var req *http.Request
	if req, err = c.prepareRequest(options); err != nil {
		return
	}

	if resp, err = c.doIt(req, options.Timeout); err != nil {
		return
	}

	switch resp.StatusCode {
	case 200:
	default:
		return nil, errors.New(resp.Status)
	}

	return
}

// Download ...
func (c *crawler) Download(options web.Request) (filePath string, err error) {

	log.Infof("Downloading %s", options.Url)

	// get heads
	var headResp *http.Response
	if headResp, err = c.HEAD(options); err != nil {
		return
	}

	defer headResp.Body.Close()
	size, _ := strconv.ParseInt(headResp.Header.Get("Content-Length"), 10, 64)

	log.Infof("File %s size: %s", options.Url, prettyByteSize(size))

	var req *http.Request
	if req, err = c.prepareRequest(options); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = c.doIt(req, options.Timeout); err != nil {
		return
	}

	defer resp.Body.Close()

	var f *os.File
	f, err = os.CreateTemp("/tmp", "*")
	if err != nil {
		return
	}
	defer f.Close()

	done := make(chan struct{})
	defer close(done)

	go printDownloadPercent(done, options.Url, f, size)

	filePath = f.Name()
	if _, err = io.Copy(f, resp.Body); err != nil {
		return
	}

	log.Infof("File %s downloaded to: %s", options.Url, filePath)

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

func (c *crawler) prepareRequest(options web.Request) (req *http.Request, err error) {

	if options.Url == "" {
		err = apperr.ErrBadRequestParams
		return
	}

	var r io.Reader = nil
	switch options.Method {
	case "POST", "PUT", "UPDATE", "PATCH", "HEAD":
		if options.Body != nil && len(options.Body) > 0 {
			r = bytes.NewReader(options.Body)
		}
	}

	if options.Context == nil {
		options.Context = context.Background()
	}

	if req, err = http.NewRequestWithContext(options.Context, options.Method, options.Url, r); err != nil {
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15")

	// set headers
	for _, values := range options.Headers {
		for k, v := range values {
			req.Header.Set(k, v)
		}
	}

	if err = c.auth(req, options.Url); err != nil {
		return
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

func printDownloadPercent(done chan struct{}, uri string, file *os.File, total int64) {
	startTime := time.Now()
	totalPretty := prettyByteSize(total)
	var oldValue = 0.0
	for {
		now := time.Now()
		select {
		case <-done:
			return
		default:
			fi, err := file.Stat()
			if err != nil {
				log.Error(err.Error())
			}
			size := fi.Size()
			if size == 0 {
				size = 1
			}
			var percent = float64(size) / float64(total) * 100
			showPercent := math.Floor(percent)
			if showPercent == oldValue {
				continue
			}
			currentPretty := prettyByteSize(size)
			eta := now.Sub(startTime) / time.Duration(showPercent) * time.Duration(100-showPercent)
			log.Infof("Downloading %s: %s/%s (%.2f%%) ETA: %s", uri, currentPretty, totalPretty, showPercent, eta)
			oldValue = showPercent
		}
		time.Sleep(time.Second / 2)
	}
}

func prettyByteSize(b int64) string {
	bf := float64(b)
	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sB", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", bf)
}
