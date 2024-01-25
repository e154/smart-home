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

package uptime

import (
	"context"
	"math"
	"time"

	"go.uber.org/atomic"

	"github.com/e154/smart-home/common"
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/system/supervisor"
)

// Actor ...
type Actor struct {
	*supervisor.BaseActor
	appStarted time.Time
	total      *atomic.Uint64
}

// NewActor ...
func NewActor(entity *m.Entity,
	service supervisor.Service) *Actor {

	attributes := NewAttr()
	attributes.Deserialize(entity.Attributes.Serialize())
	entity.Attributes = attributes

	actor := &Actor{
		BaseActor:  supervisor.NewBaseActor(entity, service),
		appStarted: time.Now(),
		total:      atomic.NewUint64(0),
	}
	if actor.Attrs == nil {
		actor.Attrs = NewAttr()
	}
	return actor
}

func (e *Actor) Destroy() {

}

func (e *Actor) Spawn() {
	go func() {
		e.check()
		e.update()
	}()
}

func (e *Actor) update() {

	total, err := GetUptime()
	if err != nil {
		return
	}

	e.total.Store(total)

	e.AttrMu.Lock()
	e.Attrs[AttrUptimeTotal].Value = e.total.Load()
	e.AttrMu.Unlock()

	e.SaveState(false, false)
}

func (e *Actor) check() {

	var lastStart time.Time
	var lastShutdown *time.Time

	ctx := context.Background()
	history, _, err := e.Service.Adaptors().RunHistory.List(ctx, 1, 1, "desc", "start", nil)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if len(history) > 0 {
		lastShutdown = history[0].End
	}

	attributes := NewAttr()
	var from *time.Time
	if f, ok := attributes[AttrLastStart]; ok {
		if f.Time().String() != "0001-01-01 00:00:00 +0000 UTC" {
			from = common.Time(f.Time())
		}
	}

	var uptime = attributes[AttrUptime].Int64()
	var downtime int64

	var firstStart *time.Time
	if f, ok := attributes[AttrFirstStart]; ok {
		if f.Time().String() != "0001-01-01 00:00:00 +0000 UTC" {
			firstStart = common.Time(f.Time())
		}
	}

	const perPage int64 = 500
	var page int64 = 0
LOOP:

	history, _, err = e.Service.Adaptors().RunHistory.List(ctx, perPage, page*perPage, "asc", "start", from)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if page == 0 && len(history) > 0 && firstStart == nil {
		firstStart = common.Time(history[0].Start)
	}

	for _, story := range history {
		if story.End != nil {
			uptime += int64(story.End.Sub(story.Start).Seconds())
		}
	}

	if len(history) != 0 {
		lastStart = history[len(history)-1].Start
		page++
		goto LOOP
	}

	now := time.Now()
	timeFromFirstStart := int64(now.Sub(*firstStart).Seconds())

	downtime = timeFromFirstStart - uptime

	uptimePercent := float64(uptime) / float64(timeFromFirstStart) * 100
	uptimePercent = math.Round(uptimePercent*100) / 100
	downtimePercent := float64(downtime) / float64(timeFromFirstStart) * 100
	downtimePercent = math.Round(downtimePercent*100) / 100

	if uptime == 0 {
		uptimePercent = 100
		downtimePercent = 0
	}

	var attributeValues = make(m.AttributeValue)
	attributeValues[AttrAppStarted] = e.appStarted
	attributeValues[AttrFirstStart] = common.TimeValue(firstStart)
	if lastShutdown != nil {
		attributeValues[AttrLastShutdown] = common.TimeValue(lastShutdown)
	}
	attributeValues[AttrLastShutdownCorrect] = lastShutdown != nil
	attributeValues[AttrLastStart] = lastStart
	attributeValues[AttrUptime] = uptime
	attributeValues[AttrDowntime] = downtime
	attributeValues[AttrUptimePercent] = uptimePercent
	attributeValues[AttrDowntimePercent] = downtimePercent

	e.DeserializeAttr(attributeValues)

	e.SaveState(false, true)

	//fmt.Println("first start", firstStart)
	//fmt.Println("last shutdown", lastShutdown)
	//fmt.Println("last shutdown correct", lastShutdown != nil)
	//fmt.Println("last start", lastStart)
	//fmt.Println("now", now)
	//fmt.Println("uptime", uptime, uptimePercent)
	//fmt.Println("downtime", downtime, downtimePercent)
}
