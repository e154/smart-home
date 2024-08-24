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

package adaptors

import (
	"context"

	"gorm.io/gorm"

	"github.com/e154/smart-home/db"
	m "github.com/e154/smart-home/models"
)

// IAutomation ...
type IAutomation interface {
	Statistic(ctx context.Context) (statistic *m.AutomationStatistic, err error)
}

// Automation ...
type Automation struct {
	IAutomation
	table *db.Automation
	db    *gorm.DB
}

// GetAutomationAdaptor ...
func GetAutomationAdaptor(d *gorm.DB) IAutomation {
	return &Automation{
		table: &db.Automation{Db: d},
		db:    d,
	}
}

func (n *Automation) Statistic(ctx context.Context) (statistic *m.AutomationStatistic, err error) {
	var dbVer *db.AutomationStatistic
	if dbVer, err = n.table.Statistic(ctx); err != nil {
		return
	}
	statistic = &m.AutomationStatistic{
		TasksTotal:      dbVer.TasksTotal,
		TasksEnabled:    dbVer.TasksEnabled,
		TriggersTotal:   dbVer.TriggersTotal,
		TriggersEnabled: dbVer.TriggersEnabled,
		ConditionsTotal: dbVer.ConditionsTotal,
		ActionsTotal:    dbVer.ActionsTotal,
	}
	return
}
