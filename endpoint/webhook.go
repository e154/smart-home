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

package endpoint

import (
	"github.com/e154/smart-home/common/apperr"
	"github.com/pkg/errors"
	"net/http"
)

// WebhookEndpoint ...
type WebhookEndpoint struct {
	*CommonEndpoint
}

// NewWebhookEndpoint ...
func NewWebhookEndpoint(common *CommonEndpoint) *WebhookEndpoint {
	return &WebhookEndpoint{
		CommonEndpoint: common,
	}
}

// Webhook ...
func (p *WebhookEndpoint) Webhook(w http.ResponseWriter, r *http.Request) (err error) {

	if isLoaded := p.supervisor.PluginIsLoaded("webhook"); !isLoaded {
		err = errors.Wrap(apperr.ErrInternal, "plugin not loaded")
		return
	}

	var pl interface{}
	if pl, err = p.supervisor.GetPlugin("webhook"); err != nil {
		return
	}

	plugin, ok := pl.(http.Handler)
	if !ok {
		return
	}

	plugin.ServeHTTP(w, r)

	return
}
