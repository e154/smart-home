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

package alexa

import (
	"errors"
	"github.com/e154/smart-home/system/cache"
)

type AppSession struct {
	pull map[string]cache.Cache
}

func NewAppSession() *AppSession {
	return &AppSession{
		pull: make(map[string]cache.Cache),
	}
}

func (h *AppSession) addSession(session string) (c cache.Cache, err error) {
	if c, err = cache.NewCache("memory", `{"interval":3600}`); err != nil {
		return
	}
	h.pull[session] = c
	return
}

func (h *AppSession) getSession(session string) (c cache.Cache, err error) {
	var exist bool
	if c, exist = h.pull[session]; !exist {
		err = errors.New("record not found")
		return
	}
	return
}

func (h *AppSession) delSession(session string) {
	if _, exist := h.pull[session]; exist {
		delete(h.pull, session)
	}
	return
}
