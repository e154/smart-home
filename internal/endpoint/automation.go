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
	"context"

	"github.com/e154/smart-home/pkg/models"
)

// AutomationEndpoint ...
type AutomationEndpoint struct {
	*CommonEndpoint
}

// NewAutomationEndpoint ...
func NewAutomationEndpoint(common *CommonEndpoint) *AutomationEndpoint {
	return &AutomationEndpoint{
		CommonEndpoint: common,
	}
}

// Statistic ...
func (n *AutomationEndpoint) Statistic(ctx context.Context) (statistic []*models.Statistic, err error) {
	var stat *models.AutomationStatistic
	if stat, err = n.adaptors.Automation.Statistic(ctx); err != nil {
		return
	}
	statistic = []*models.Statistic{
		{
			Name:        "automation.stat_tasks_total_name",
			Description: "automation.stat_tasks_total_descr",
			Value:       stat.TasksTotal,
		},
		{
			Name:        "automation.stat_tasks_enabled_name",
			Description: "automation.stat_tasks_enabled_descr",
			Value:       stat.TasksEnabled,
		},
		{
			Name:        "automation.stat_triggers_total_name",
			Description: "automation.stat_triggers_total_descr",
			Value:       stat.TriggersTotal,
		},
		{
			Name:        "automation.stat_triggers_enabled_name",
			Description: "automation.stat_triggers_enabled_descr",
			Value:       stat.TriggersEnabled,
		},
		{
			Name:        "automation.stat_conditions_total_name",
			Description: "automation.stat_conditions_total_descr",
			Value:       stat.ConditionsTotal,
		},
		{
			Name:        "automation.stat_actions_total_name",
			Description: "automation.stat_actions_total_descr",
			Value:       stat.ActionsTotal,
		},
	}
	return
}
