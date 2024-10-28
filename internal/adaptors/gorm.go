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

	"github.com/e154/smart-home/internal/system/migrations"
	"github.com/e154/smart-home/internal/system/orm"
	"github.com/e154/smart-home/pkg/adaptors"
	"github.com/e154/smart-home/pkg/logger"
	"github.com/e154/smart-home/pkg/models"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var (
	log = logger.MustGetLogger("adaptors")
)

type Adaptors struct {
	*adaptors.Adaptors
}

// NewAdaptors ...
func NewAdaptors(lc fx.Lifecycle,
	db *gorm.DB,
	cfg *models.AppConfig,
	migrations *migrations.Migrations,
	orm *orm.Orm) (a *adaptors.Adaptors) {

	a = &adaptors.Adaptors{
		Transaction:       NewTransactionManger(db),
		Script:            GetScriptAdaptor(db),
		Tag:               GetTagAdaptor(db),
		Role:              GetRoleAdaptor(db),
		Permission:        GetPermissionAdaptor(db),
		User:              GetUserAdaptor(db),
		UserMeta:          GetUserMetaAdaptor(db),
		UserDevice:        GetUserDeviceAdaptor(db),
		Image:             GetImageAdaptor(db),
		Variable:          GetVariableAdaptor(db),
		Entity:            GetEntityAdaptor(db, orm),
		EntityState:       GetEntityStateAdaptor(db),
		EntityAction:      GetEntityActionAdaptor(db),
		EntityStorage:     GetEntityStorageAdaptor(db),
		Log:               GetLogAdaptor(db),
		Template:          GetTemplateAdaptor(db),
		Message:           GetMessageAdaptor(db),
		MessageDelivery:   GetMessageDeliveryAdaptor(db),
		Zigbee2mqtt:       GetZigbee2mqttAdaptor(db),
		Zigbee2mqttDevice: GetZigbee2mqttDeviceAdaptor(db),
		AlexaSkill:        GetAlexaSkillAdaptor(db),
		AlexaIntent:       GetAlexaIntentAdaptor(db),
		Metric:            GetMetricAdaptor(db, orm),
		MetricBucket:      GetMetricBucketAdaptor(db, orm),
		Area:              GetAreaAdaptor(db),
		Action:            GetActionAdaptor(db, orm),
		Condition:         GetConditionAdaptor(db),
		Trigger:           GetTriggerAdaptor(db, orm),
		Task:              GetTaskAdaptor(db, orm),
		RunHistory:        GetRunHistoryAdaptor(db),
		Plugin:            GetPluginAdaptor(db),
		TelegramChat:      GetTelegramChannelAdaptor(db),
		Dashboard:         GetDashboardAdaptor(db),
		DashboardTab:      GetDashboardTabAdaptor(db),
		DashboardCard:     GetDashboardCardAdaptor(db),
		DashboardCardItem: GetDashboardCardItemAdaptor(db),
		ScriptVersion:     GetScriptVersionAdaptor(db),
		Automation:        GetAutomationAdaptor(db),
	}

	if lc != nil {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				if cfg != nil && migrations != nil && cfg.AutoMigrate {
					err = migrations.Up()
					return
				}
				return
			},
		})
	}

	return
}
