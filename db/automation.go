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

package db

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/e154/smart-home/common/apperr"
)

// Automation ...
type Automation struct {
	Db *gorm.DB
}

type AutomationStatistic struct {
	TasksTotal      int32 `gorm:"->"`
	TasksEnabled    int32 `gorm:"->"`
	TriggersTotal   int32 `gorm:"->"`
	TriggersEnabled int32 `gorm:"->"`
	ConditionsTotal int32 `gorm:"->"`
	ActionsTotal    int32 `gorm:"->"`
}

// Statistic ...
func (n *Automation) Statistic(ctx context.Context) (statistic *AutomationStatistic, err error) {

	statistic = &AutomationStatistic{}

	err = n.Db.WithContext(ctx).Raw(`
select tasks_total.count as tasks_total, tasks_enabled.count as tasks_enabled,
       triggers_total.count as triggers_total, triggers_enabled.count as triggers_enabled,
       conditions_total.count as conditions_total, actions_total.count as actions_total
from (select count(id)
      from tasks) as tasks_total,
     (select count(id)
      from tasks
      where enabled = true) as tasks_enabled,

     (select count(id)
      from triggers) as triggers_total,
     (select count(id)
      from triggers
      where enabled = true) as triggers_enabled,

     (select count(id)
      from conditions) as conditions_total,

     (select count(id)
      from actions) as actions_total;`).
		Scan(statistic).
		Error

	if err != nil {
		err = errors.Wrap(apperr.ErrAutomationStat, err.Error())
	}

	return
}
