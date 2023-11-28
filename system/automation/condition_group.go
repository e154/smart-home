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

package automation

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
)

// NewConditionGroup ...
func NewConditionGroup(t common.ConditionType) *ConditionGroup {
	return &ConditionGroup{
		t:          t,
		lastStatus: atomic.NewBool(false),
		Mutex:      sync.Mutex{},
	}
}

// ConditionGroup ...
type ConditionGroup struct {
	rules      []*Condition
	t          common.ConditionType
	lastStatus *atomic.Bool
	sync.Mutex
}

// AddCondition ...
func (c *ConditionGroup) AddCondition(condition *Condition) {
	c.rules = append(c.rules, condition)
}

func (c *ConditionGroup) Stop() {
	for _, condition := range c.rules {
		condition.Stop()
	}
}

// Check ...
func (c *ConditionGroup) Check(entityId *common.EntityId) (state bool, err error) {
	c.Lock()
	defer func() {
		c.lastStatus.Store(state)
		c.Unlock()
	}()

	var total = len(c.rules)
	if total == 0 {
		state = true
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(total)

	for _, r := range c.rules {
		go func(condition *Condition) {
			defer wg.Done()
			bg, _ := context.WithTimeout(context.Background(), time.Second)
			ctx := context.WithValue(bg, "entityId", entityId)
			if _, err = condition.Check(ctx); err != nil {
				log.Error(err.Error())
			}
		}(r)
	}

	wg.Wait()

	switch c.t {
	case common.ConditionOr:
		for _, r := range c.rules {
			if r.Status() {
				state = true
				return
			}
		}

	case common.ConditionAnd:
		state = true
		for _, r := range c.rules {
			if !r.Status() {
				state = false
			}
			if !state {
				return
			}
		}
	}

	return
}
