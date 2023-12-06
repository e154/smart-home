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
	"net/http"
	"net/url"
	"regexp"
)

// HTTPRequest is a serializable version of http.Request ( with only usefull fields )
type HTTPRequest struct {
	Method        string
	URL           string
	Header        map[string][]string
	ContentLength int64
	WS            bool
}

// SerializeHTTPRequest create a new HTTPRequest from a http.Request
func SerializeHTTPRequest(req *http.Request) (r *HTTPRequest) {
	r = &HTTPRequest{
		URL:           req.URL.String(),
		Method:        req.Method,
		Header:        req.Header,
		ContentLength: req.ContentLength,
	}
	return
}

// UnserializeHTTPRequest create a new http.Request from a HTTPRequest
func UnserializeHTTPRequest(req *HTTPRequest) (r *http.Request, err error) {
	r = new(http.Request)
	r.Method = req.Method
	r.URL, err = url.Parse(req.URL)
	if err != nil {
		return
	}
	r.Header = req.Header
	r.ContentLength = req.ContentLength
	return
}

// Rule match HTTP requests to allow / deny access
type Rule struct {
	Method  string
	URL     string
	Headers map[string]string

	methodRegex  *regexp.Regexp
	urlRegex     *regexp.Regexp
	headersRegex map[string]*regexp.Regexp
}

// NewRule creates a new Rule
func NewRule(method string, url string, headers map[string]string) (rule *Rule, err error) {
	rule = new(Rule)
	rule.Method = method
	rule.URL = url
	if headers != nil {
		rule.Headers = headers
	} else {
		rule.Headers = make(map[string]string)
	}
	err = rule.Compile()
	return
}

// Compile the regular expressions
func (rule *Rule) Compile() (err error) {
	if rule.Method != "" {
		rule.methodRegex, err = regexp.Compile(rule.Method)
		if err != nil {
			return
		}
	}
	if rule.URL != "" {
		rule.urlRegex, err = regexp.Compile(rule.URL)
		if err != nil {
			return
		}
	}
	rule.headersRegex = make(map[string]*regexp.Regexp)
	for header, regexStr := range rule.Headers {
		var regex *regexp.Regexp
		regex, err = regexp.Compile(regexStr)
		if err != nil {
			return
		}
		rule.headersRegex[header] = regex
	}
	return
}

// Match returns true if the http.Request matches the Rule
func (rule *Rule) Match(req *http.Request) bool {
	if rule.methodRegex != nil && !rule.methodRegex.MatchString(req.Method) {
		return false
	}
	if rule.urlRegex != nil && !rule.urlRegex.MatchString(req.URL.String()) {
		return false
	}

	for headerName, regex := range rule.headersRegex {
		if !regex.MatchString(req.Header.Get(headerName)) {
			return false
		}
	}

	return true
}

func (rule *Rule) String() string {
	return fmt.Sprintf("%s %s %v", rule.Method, rule.URL, rule.Headers)
}
