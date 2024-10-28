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

package webhook

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/e154/smart-home/internal/common/web/urlpath"
	"github.com/e154/smart-home/internal/system/supervisor"
	m "github.com/e154/smart-home/pkg/models"
	"github.com/e154/smart-home/pkg/plugins"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	AccessToken string
	urlPath     urlpath.Path
}

// NewActor ...
func NewActor(entity *m.Entity,
	service plugins.Service) *Actor {

	settings := NewSettings()
	_, _ = settings.Deserialize(entity.Settings.Serialize())

	actor := &Actor{
		BaseActor:   supervisor.NewBaseActor(entity, service),
		AccessToken: settings[AttrToken].Decrypt(),
		urlPath:     urlpath.New(settings[AttrPath].String()),
	}

	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}

	if actor.Setts == nil {
		actor.Setts = NewSettings()
	}

	return actor
}

func (e *Actor) Destroy() {

}

// Spawn ...
func (e *Actor) Spawn() {

}

func (e *Actor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !e.checkToken(getAccessToken(r)) {
		w.WriteHeader(401)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}

	match, ok := e.urlPath.Match(r.URL.Path)
	if !ok {
		return
	}

	e.updateState(r, match)

	w.WriteHeader(200)
	w.Write([]byte("ok"))
}

func (e *Actor) checkToken(token string) bool {
	if e.AccessToken == "" {
		return true
	}
	return fmt.Sprintf("Bearer %s", e.AccessToken) == token
}

func (e *Actor) updateState(r *http.Request, match urlpath.Match) {

	var attributeValues = make(m.AttributeValue)
	var headers = []string{}
	for k, v := range r.Header {
		value := fmt.Sprintf("%v:", k)
		for i, v := range v {
			if i > 0 {
				value += ","
			}
			value += v
		}
		headers = append(headers, value)
		sort.Strings(headers)
	}
	if r != nil {
		body, _ := io.ReadAll(r.Body)
		attributeValues[AttrUrl] = r.URL.String()
		attributeValues[AttrHost] = r.Host
		attributeValues[AttrSize] = r.ContentLength
		attributeValues[AttrBody] = string(body)
		attributeValues[AttrMethod] = r.Method
		attributeValues[AttrHeaders] = strings.Join(headers[:], ";")
	}

	// save match
	for k, v := range match.Params {
		attributeValues[k] = v
	}

	// save query params
	query := r.URL.Query()
	for k, v := range query {
		attributeValues[k] = v
	}

	e.DeserializeAttr(attributeValues)

	e.SaveState(false, r != nil)
}

func getAccessToken(r *http.Request) (accessToken string) {
	accessToken = r.Header.Get("authorization")
	if accessToken != "" {
		return
	}
	accessToken = r.URL.Query().Get("access_token")
	return
}
