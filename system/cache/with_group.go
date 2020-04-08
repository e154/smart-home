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

package cache

import (
	"fmt"
	"github.com/e154/smart-home/common"
	"time"
)

var (
	log = common.MustGetLogger("cache")
)

// WithGroup ...
type WithGroup struct {
	bm        Cache
	cacheTime time.Duration
	verbose   bool
	name      string
}

// NewWithGroup ...
func NewWithGroup(name string, t time.Duration, v bool) (group *WithGroup, err error) {

	var bm Cache
	if bm, err = NewCache("memory", fmt.Sprintf(`{"interval":%d}`, t)); err != nil {
		return
	}

	group = &WithGroup{
		name:      name,
		bm:        bm,
		cacheTime: t,
		verbose:   v,
	}

	return
}

// ClearAll ...
func (c *WithGroup) ClearAll() (*WithGroup, error) {
	c.log("clear all")

	err := c.bm.ClearAll()

	return c, err
}

// GetKey ...
func (c *WithGroup) GetKey(key interface{}) string {
	return fmt.Sprintf("%s_%s", c.name, key.(string))
}

// Clear ...
func (c *WithGroup) Clear(key string) (*WithGroup, error) {
	cacheKey := c.GetKey(key)

	c.log("clear %s", cacheKey)

	err := c.bm.Delete(cacheKey)

	return c, err
}

func (c *WithGroup) addToGroup(group, key string) (*WithGroup, error) {

	g := []string{}
	w := c.bm.Get(group)
	if w != nil {
		g = w.([]string)
	}

	exist := false
	for _, v := range g {
		if key == v {
			exist = true
		}
	}

	var err error
	if !exist {
		c.log("add to group %s", group)
		g = append(g, key)
		err = c.bm.Put(group, g, c.cacheTime)
	}

	return c, err
}

// ClearGroup ...
func (c *WithGroup) ClearGroup(group string) (*WithGroup, error) {
	c.log("clear group %s", group)

	g := []string{}
	w := c.bm.Get(group)
	if w == nil {
		return c, nil
	}

	g = w.([]string)
	if len(g) == 0 {
		return c, nil
	}

	for _, key := range g {
		c.bm.Delete(key)
	}

	_, err := c.Clear(group)

	return c, err
}

// Put ...
func (c *WithGroup) Put(group, key string, val interface{}) (*WithGroup, error) {
	c.log("put key %s", key)

	if err := c.bm.Put(key, val, c.cacheTime); err != nil {
		return c, err
	}

	return c.addToGroup(group, key)
}

// IsExist ...
func (c *WithGroup) IsExist(key string) bool {

	return c.bm.IsExist(key)
}

// Get ...
func (c *WithGroup) Get(key string) interface{} {
	c.log("get key %s", key)

	return c.bm.Get(key)
}

// Delete ...
func (c *WithGroup) Delete(key string) *WithGroup {
	c.log("delete value by key %s", key)

	c.bm.Delete(key)

	return c
}

func (c *WithGroup) log(format string, a ...interface{}) {
	if c.verbose {
		log.Debugf(format, a...)
	}
}
