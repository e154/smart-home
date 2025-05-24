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

	db "github.com/e154/smart-home/internal/db"
	"github.com/e154/smart-home/pkg/adaptors"
	m "github.com/e154/smart-home/pkg/models"

	"gorm.io/gorm"
)

var _ adaptors.AutomationRepo = (*Automation)(nil)

// Automation ...
type Automation struct {
	table *db.Automation
	db    *gorm.DB
}

// GetAutomationAdaptor ...
func GetAutomationAdaptor(d *gorm.DB) *Automation {
	return &Automation{
		table: &db.Automation{&db.Common{Db: d}},
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
